 <template>
  <div id="container">
    <div class="a"></div>
    <div class="b"></div>
    <div class="c"></div>
    <div class="d"></div>
    <div class="e">
    <img :src="normal_gt_img"></div>
    <div class="f">
    <img :src="color_img"></div>
  </div>
</template>
    
    <script>
import * as three from "three";
import { GUI } from "dat.gui";
import "three-orbitcontrols";
import { Vector3 } from 'three';
export default {
  name: "Threejs",
  data() {
    return {
      camera: null,
      scene: null,
      renderer: null,
      mesh: null,
      height: null,
      width: null,
      gui: null,
      ori_position: new Vector3(0,0,0),
      AMOUNT: 4,
      cameras: [],
      renderers: [],
      scenes: [],
      controls: [],
      plylist: ["mvsnet", "rmvsnet", "var","low_gt"],
      id_list: [...Array(90).keys()],
      vertex_shader: {
        depth: `
        attribute float alpha;
        attribute vec3 color;
        varying vec3 vcolor;
        uniform float mn;
        uniform float mx;
        uniform vec3 pos;
        void main() {
          vec4 mvPosition = modelViewMatrix * vec4( position, 1.0 );
          gl_PointSize = 2.0;
          gl_Position = projectionMatrix * mvPosition;
          if(position.x>mx) vcolor = vec3(0,0,0);
          else vcolor = vec3((position.x-mn)/(mx-mn),(position.x-mn)/(mx-mn),(position.x-mn)/(mx-mn));
        }`,
        prob: `
        attribute float prob;
        attribute float alpha;
        attribute vec3 color;
        varying vec3 vcolor;
        uniform float mn;
        uniform float mx;
        void main() {
          vec4 mvPosition = modelViewMatrix * vec4( position.x,position.y,position.z, 1.0 );
          gl_PointSize = 2.0;
          gl_Position = projectionMatrix * mvPosition;
          vcolor = vec3(prob,prob,prob);
        }
        `,
        normal: `
        attribute float prob;
        attribute float alpha;
        varying vec3 vcolor;
        uniform float mn;
        uniform float mx;
        void main() {
          vec4 mvPosition = modelViewMatrix * vec4( position.x,position.y,position.z, 1.0 );
          gl_PointSize = 2.0;
          gl_Position = projectionMatrix * mvPosition;
          vcolor = vec3((normal.x+1.0)/2.0,(normal.y+1.0)/2.0,(normal.z+1.0)/2.0);
        }
        `,
        white: `
        attribute float prob;
        attribute float alpha;
        varying vec3 vcolor;
        uniform float mn;
        uniform float mx;
        void main() {
          vec4 mvPosition = modelViewMatrix * vec4( position.x,position.y,position.z, 1.0 );
          gl_PointSize = 2.0;
          gl_Position = projectionMatrix * mvPosition;
          vcolor = vec3(1,1,1);
        }
        `,
        err: `
        attribute float prob;
        attribute float err;
        attribute float alpha;
        varying vec3 vcolor;
        uniform float mn;
        uniform float mx;
        float colormap_red(float x) {
            if (x < 0.7) {
                return 4.0 * x - 1.5;
            } else {
                return -4.0 * x + 4.5;
            }
        }

        float colormap_green(float x) {
            if (x < 0.5) {
                return 4.0 * x - 0.5;
            } else {
                return -4.0 * x + 3.5;
            }
        }

        float colormap_blue(float x) {
            if (x < 0.3) {
              return 4.0 * x + 0.5;
            } else {
              return -4.0 * x + 2.5;
            }
        }

        vec3 colormap(float x) {
            float r = clamp(colormap_red(x), 0.0, 1.0);
            float g = clamp(colormap_green(x), 0.0, 1.0);
            float b = clamp(colormap_blue(x), 0.0, 1.0);
            return vec3(r, g, b);
        }
        void main() {
          vec4 mvPosition = modelViewMatrix * vec4( position.x,position.y,position.z, 1.0 );
          gl_PointSize = 2.0;
          gl_Position = projectionMatrix * mvPosition;
          if(err>mx) {vcolor = vec3(1.0,1.0,1.0);}
          else {vcolor = colormap(err/mx);}
        }
        `,
      },
      fragment_shader: `
        varying vec3 vcolor;
        void main() {
          gl_FragColor = vec4( vcolor,1 );
        }`,
      Controller: {
        x: 0,
        y: 0,
        z: 0,
        rx: 0,
        ry: 1.57,
        rz: 1.57,
        depthR: 10000,
        shader: "white",
        err_range:128,
        cur_id: 0
      },
      color_img:"",
      normal_gt_img:""
    };
  },
  computed: {
  },
  methods: {
    async setScenes() {
      for (var i = 0; i < this.AMOUNT; i++) {
        var res = await this.$axios.get("file/"+this.id_list[this.Controller.cur_id]+"/"+this.plylist[i]);
        var points = new Float32Array(res.data.points);
        var probs = new Float32Array(res.data.probs);
        var errs = new Float32Array(res.data.errs);
        var geo = new three.BufferGeometry();
        geo.setAttribute("position", new three.BufferAttribute(points, 3));
        geo.setAttribute("prob", new three.BufferAttribute(probs, 1));
        geo.setAttribute("err", new three.BufferAttribute(errs, 1));
        var uniforms = {
          mn: { value: 0 },
          mx: { value: 5000 },
          pos: {value: this.ori_position}
        };
        var shaderMaterial = new three.ShaderMaterial({
          uniforms: uniforms,
          vertexShader: this.vertex_shader[this.Controller.shader],
          fragmentShader: this.fragment_shader,
          transparent: true,
        });
        var cloud = new three.Points(geo, shaderMaterial);
        cloud.name = "points";
        cloud.geometry.computeVertexNormals();
        var scene = new three.Scene();
        scene.add(new three.AmbientLight(0x222244));
        scene.add(cloud);
        if(this.scenes.length>i){
          this.scenes[i] = scene;
        }else{
          this.scenes.push(scene);
        }
        
        if(this.renderers.length<=i){
          var renderer = new three.WebGLRenderer();
          renderer.setPixelRatio(window.devicePixelRatio);
          renderer.setSize(this.width, this.height);
          renderer.shadowMap.enabled = true;
          this.renderers.push(renderer);
            document.body
            .getElementsByClassName(String.fromCharCode(97 + i))[0]
          .appendChild(renderer.domElement);
        }
      }
    },
    async setControllers() {
      var res = await this.$axios.get("calib/"+this.id_list[this.Controller.cur_id]);
      var camera_pos = new Float32Array(res.data.camera_pos);
      this.camera = new three.PerspectiveCamera(
        90,
        this.width / this.height,
        0.1,
        10000
      );
      var m = new three.Matrix4();
      console.log("camera pos is",camera_pos);
      m.set(camera_pos[0], camera_pos[1], camera_pos[2], camera_pos[3], camera_pos[4], camera_pos[5], camera_pos[6], camera_pos[7], camera_pos[8], camera_pos[9], camera_pos[10], camera_pos[11], camera_pos[12], camera_pos[13], camera_pos[14], camera_pos[15]);
      console.log("m is ",m);
      m.decompose(
        this.camera.position,
        this.camera.quaternion,
        this.camera.scale
      );
      var t = this.camera.position.x;
      this.camera.position.x = this.camera.position.y;
      this.camera.position.y = t;
      this.camera.position.z = - this.camera.position.z;
      t = this.camera.rotation.x;
      this.camera.rotation.x = this.camera.rotation.y;
      this.camera.rotation.y = -t;
      this.camera.rotation.z = - this.camera.rotation.z;
      console.log(this.camera.position,this.camera.rotation,this.camera.quaternion)
      var cam = new three.PerspectiveCamera(
        90,
        this.width / this.height,
        0.1,
        10000
      );
      for (var i = 0; i < this.AMOUNT; i++) {
        var control = new three.OrbitControls(
          cam,
          this.renderers[i].domElement
        );
        control.addEventListener("change", this.render); // call this only in static scenes (i.e., if there is no animation loop)
        control.enabled = true;
        this.controls.push(control);
      }
      this.ori_position = this.camera.position;
      for(i=0;i<this.AMOUNT;i++){
        this.controls[i].update();
      }
      this.Controller.x = this.camera.position.x;
      this.Controller.y = this.camera.position.y;
      this.Controller.z = this.camera.position.z;
      this.Controller.rx = this.camera.rotation.x;
      this.Controller.ry = this.camera.rotation.y;
      this.Controller.rz = this.camera.rotation.z;
    },
    async setGUI() {
      this.gui = new GUI();
      var tmp;
      tmp = this.gui.add(this.Controller, "x", -10000, 10000);
      tmp.onChange((value) => {
        this.camera.position.x = value;
        for(var i=0;i<this.AMOUNT;i++) 
          this.controls[i].update();
        this.render();
      });
      tmp.listen();
      tmp = this.gui.add(this.Controller, "y", -10000, 10000);
      tmp.listen();
      tmp.onChange((value) => {
        this.camera.position.y = value;
        for(var i=0;i<this.AMOUNT;i++) 
          this.controls[i].update();
        this.render();
      });
      tmp = this.gui.add(this.Controller, "z", -10000, 10000);
      tmp.listen();
      tmp.onChange((value) => {
        this.camera.position.z = value;
        for(var i=0;i<this.AMOUNT;i++) 
          this.controls[i].update();
        this.render();
      });
      tmp = this.gui.add(this.Controller, "rx", -3.14,3.14);
      tmp.listen();
      tmp.onChange((value) => {
        this.camera.rotation.x = value;
        for(var i=0;i<this.AMOUNT;i++) 
          this.controls[i].update();
        this.render();
      });
      tmp = this.gui.add(this.Controller, "ry", -3.14,3.14);
      tmp.listen();
      tmp.onChange((value) => {
        this.camera.rotation.y = value;
        for(var i=0;i<this.AMOUNT;i++) 
          this.controls[i].update();
        this.render();
      });
      tmp = this.gui.add(this.Controller, "rz", -3.14,3.14);
      tmp.listen();
      tmp.onChange((value) => {
        this.camera.rotation.z = value;
        for(var i=0;i<this.AMOUNT;i++) 
          this.controls[i].update();
        this.render();
      });
      tmp = this.gui.add(this.Controller, "depthR", 0, 10000);
      tmp.listen();
      tmp.onChange((v) => {
        this.Controller.depthR = v;
        var uniforms = {
          color: { value: new three.Color(0xffffff) },
          mn: { value: 0 },
          mx: { value: v },
        };
        var shaderMaterial = new three.ShaderMaterial({
          uniforms: uniforms,
          vertexShader: this.vertex_shader[this.Controller.shader],
          fragmentShader: this.fragment_shader,
          transparent: true,
        });
        for (var i = 0; i < this.AMOUNT; i++) {
          this.scenes[i].getObjectByName("points").material = shaderMaterial;
        }
        this.render();
      });
      tmp = this.gui.add(this.Controller, "err_range", 0, 1000);
      tmp.listen();
      tmp.onChange((v) => {
        this.Controller.err_range = v;
        var uniforms = {
          color: { value: new three.Color(0xffffff) },
          mn: { value: 0 },
          mx: { value: v },
          pos: {value: new Vector3(this.ori_position)}
        };
        var shaderMaterial = new three.ShaderMaterial({
          uniforms: uniforms,
          vertexShader: this.vertex_shader[this.Controller.shader],
          fragmentShader: this.fragment_shader,
          transparent: true,
        });
        for (var i = 0; i < this.AMOUNT; i++) {
          this.scenes[i].getObjectByName("points").material = shaderMaterial;
        }
        this.render();
      });
      tmp = this.gui.add(
        this.Controller,
        "shader",
        Object.keys(this.vertex_shader)
      );
      tmp.listen();
      tmp.onChange((s) => {
        console.log(s);
        this.Controller.shader = s;
        var uniforms = {
          color: { value: new three.Color(0xffffff) },
          mn: { value: 0 },
          mx: { value: this.Controller.depthR },
          pos: {value: new Vector3(this.ori_position)}
        };
        var shaderMaterial = new three.ShaderMaterial({
          uniforms: uniforms,
          vertexShader: this.vertex_shader[this.Controller.shader],
          fragmentShader: this.fragment_shader,
          transparent: true,
        });
        for (var i = 0; i < this.AMOUNT; i++) {
          this.scenes[i].getObjectByName("points").material = shaderMaterial;
        }
        this.render();
      });
      tmp = this.gui.add(
        this.Controller,
        "cur_id",
        this.id_list
      );
      tmp.listen();
      tmp.onChange((s) => {
        this.Controller.cur_id = s;
        this.reinit();
      });
    },
    async setImg(){
      this.normal_gt_img = "img/" + this.Controller.cur_id +  "/normal_gt"
      this.color_img = "img/" + this.Controller.cur_id +  "/img_0"
    },
    async init() {
      
      this.scenes = [];
      this.renderers = [];
      this.controls = [];
      
      this.width = parseInt(
        window
          .getComputedStyle(document.body.getElementsByClassName("a")[0])
          .getPropertyValue("width"),
        10
      );
      this.height = parseInt(
        window
          .getComputedStyle(document.body.getElementsByClassName("a")[0])
          .getPropertyValue("height"),
        10
      );
      await this.setImg();
      await this.setScenes();
      await this.setControllers();
      await this.setGUI();
      this.render();
    },
    
    async reinit() {
      
      this.width = parseInt(
        window
          .getComputedStyle(document.body.getElementsByClassName("a")[0])
          .getPropertyValue("width"),
        10
      );
      this.height = parseInt(
        window
          .getComputedStyle(document.body.getElementsByClassName("a")[0])
          .getPropertyValue("height"),
        10
      );
      await this.setImg();
      await this.setScenes();
      await this.setControllers();
      this.render();
    },
    animate: function () {
      requestAnimationFrame(this.animate);
      this.mesh.rotation.x += 0.01;
      this.mesh.rotation.y += 0.02;
      this.renderer.render(this.scene, this.camera);
    },
    render: function () {
      this.Controller.x = this.camera.position.x;
      this.Controller.y = this.camera.position.y;
      this.Controller.z = this.camera.position.z;
      this.Controller.rx = this.camera.rotation.x;
      this.Controller.ry = this.camera.rotation.y;
      this.Controller.rz = this.camera.rotation.z;

      this.camera.updateProjectionMatrix();
      for (var i = 0; i < this.AMOUNT; i++) {
        this.renderers[i].render(this.scenes[i], this.camera);
      }
    },
  },
  mounted() {
    this.init();
  },
};
</script>
    
    <style scoped>
#container {
  display: grid;
  height: 100%;
  width: 100%;
  position: absolute;
  top: 0px;
  bottom: 0px;
  grid-template-columns: 30% 30% 30% 10%;
  grid-template-rows: 50% 50%;
  grid-template-areas:
    "a b c g"
    "d e f g";
}

.a {
  grid-area: a;
  background-color: aliceblue;
}

.b {
  grid-area: b;
  background-color: aqua;
}

.c {
  grid-area: c;
  background-color: blue;
}

.d {
  grid-area: d;
  background-color: brown;
}

.e {
  grid-area: e;
  background-color: red;
}

.f {
  grid-area: f;
  background-color: yellow;
}

img {
    height: 100%;
}
</style>
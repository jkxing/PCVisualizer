 <template>
  <div id="container">
    <div class="a"></div>
    <div class="b"></div>
    <div class="c"></div>
    <div class="d"></div>
    <div class="e"></div>
    <div class="f"></div>
  </div>
</template>
    
    <script>
import * as three from "three";
import { GUI } from "dat.gui";
import "three-orbitcontrols";
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
      AMOUNT: 3,
      cameras: [],
      renderers: [],
      scenes: [],
      controls: [],
      plylist: [
        "mvs.xyz",
        "rmvs.xyz",
        "var.xyz",
      ],
      vertex_shader: `
        attribute float alpha;
        attribute vec3 color;
        varying vec3 vcolor;
        uniform float mn;
        uniform float mx;
        void main() {
          vec4 mvPosition = modelViewMatrix * vec4( position, 1.0 );
          gl_PointSize = 2.0;
          gl_Position = projectionMatrix * mvPosition;
          if(position.z>mx) vcolor = vec3(0,0,0);
          else vcolor = vec3((position.z-mn)/(mx-mn),(position.z-mn)/(mx-mn),(position.z-mn)/(mx-mn));
        }`,
      fragment_shader: `
        varying vec3 vcolor;
        void main() {
          gl_FragColor = vec4( 1,1,1,1 );
        }`,
    };
  },
  methods: {
    async init() {
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
      var ASPECT_RATIO = this.width / this.height;
      for (var i = 0; i < this.AMOUNT; i++) {
        var res = await this.$axios.get(this.plylist[i]);
        //var camera_pos = res.data.camera_pos;
        var points = res.data.points;
        //var err = res.data.err;
        var geo = new three.Geometry();
        for (var j = 0; j < points.length; j += 4) {
          geo.vertices.push(
            new three.Vector3(points[j], points[j + 1], points[j + 2])
          );
        }
        var uniforms = {
          mn: { value: 0 },
          mx: { value: 5000 }
        }
        var shaderMaterial = new three.ShaderMaterial({
          uniforms: uniforms,
          vertexShader: this.vertex_shader,
          fragmentShader: this.fragment_shader,
          transparent: true
        });
        var cloud = new three.Points(geo, shaderMaterial);
        cloud.name = "points";
        var subcamera = new three.PerspectiveCamera(90, ASPECT_RATIO, 0.1, 10000);
        
        var m = new three.Matrix4();
        m.set(-0.031339678913354874,-0.999508798122406,-7.87170577609686e-08,2361.243963241577,
-0.14191657304763794, 0.0044498806819319725, -0.9898686408996582, 762.1212005615234,
0.9893823862075806, -0.03102215752005577, -0.14198632538318634, -726.7777919769287,
0.0, 0.0, 0.0, 1.0);
        m.decompose(subcamera.position,subcamera.rotation,subcamera.scale);
        console.log(subcamera.position,subcamera.rotation,subcamera.scale);
        //subcamera.rotation.set(camera_pos[3],camera_pos[4],camera_pos[5]);
        this.cameras.push(subcamera);
        var scene = new three.Scene();
        scene.add(new three.AmbientLight(0x222244));
        var geometry = new three.PlaneBufferGeometry(100, 100);
        var material = new three.MeshPhongMaterial({ color: 0x000066 });
        var background = new three.Mesh(geometry, material);
        background.receiveShadow = true;
        background.position.set(0, 0, -1);
        scene.add(background);
        scene.add(cloud);
        this.scenes.push(scene);
        var renderer = new three.WebGLRenderer();
        renderer.setPixelRatio(window.devicePixelRatio);
        renderer.setSize(this.width, this.height);
        renderer.shadowMap.enabled = true;
        this.renderers.push(renderer);
        document.body
          .getElementsByClassName(String.fromCharCode(97 + i))[0]
          .appendChild(renderer.domElement);
        var control = new three.OrbitControls(subcamera, renderer.domElement);
        control.addEventListener("change", this.render); // call this only in static scenes (i.e., if there is no animation loop)
        control.enabled = true;
        this.controls.push(control);
      }

      var effectController = {
        x: 0,
        y: 0,
        z: 0,
        rx: 0,
        ry: 0,
        rz: 0,
        depthR: 0,
      };
      this.gui = new GUI();

      this.gui.add(effectController, "x", -10, 10).onChange((value) => {
        for (var i = 0; i < this.AMOUNT; i++) {
          this.cameras[i].position.x = value;
        }
        this.render();
      });
      this.gui.add(effectController, "y", -10, 10).onChange((value) => {
        for (var i = 0; i < this.AMOUNT; i++) {
          this.cameras[i].position.y = value;
        }
        this.render();
      });
      this.gui.add(effectController, "z", -10, 10).onChange((value) => {
        for (var i = 0; i < this.AMOUNT; i++) {
          this.cameras[i].position.z = value;
        }
        this.render();
      });
      this.gui.add(effectController, "rx", -3.14, 3.14).onChange((value) => {
        for (var i = 0; i < this.AMOUNT; i++) {
          this.cameras[i].rotation.x = value;
        }
        this.render();
      });
      this.gui.add(effectController, "ry", -3.14, 3.14).onChange((value) => {
        for (var i = 0; i < this.AMOUNT; i++) {
          this.cameras[i].rotation.y = value;
        }
        this.render();
      });
      this.gui.add(effectController, "rz", -3.14, 3.14).onChange((value) => {
        for (var i = 0; i < this.AMOUNT; i++) {
          this.cameras[i].rotation.z = value;
        }
        this.render();
      });
      this.gui.add(effectController, "depthR", 0, 10).onChange((value) => {
        var uniforms = {
          color: { value: new three.Color(0xffffff) },
          mn: { value: 0 },
          mx: { value: value * 1000 },
        };

        var shaderMaterial = new three.ShaderMaterial({
          uniforms: uniforms,
          vertexShader: this.vertex_shader,
          fragmentShader: this.fragment_shader,
          transparent: true,
        });
        for (var i = 0; i < this.AMOUNT; i++) {
          this.scenes[i].getObjectByName( "points" ).material = shaderMaterial;
        }
        this.render();
      });
      this.render();

    },
    animate: function () {
      requestAnimationFrame(this.animate);
      this.mesh.rotation.x += 0.01;
      this.mesh.rotation.y += 0.02;
      this.renderer.render(this.scene, this.camera);
    },
    render: function () {
      for (var i = 0; i < this.AMOUNT; i++) {
        this.cameras[i].updateProjectionMatrix();
        this.renderers[i].render(this.scenes[i], this.cameras[i]);
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
</style>
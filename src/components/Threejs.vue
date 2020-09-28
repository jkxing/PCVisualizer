 <template>
  <div id="container"></div>
</template>
    
    <script>
import * as three from "three";
export default {
  name: "Threejs",
  data() {
    return {
      camera: null,
      scene: null,
      renderer: null,
      mesh: null,
    };
  },
  methods: {
    init: function () {
      this.camera = new three.PerspectiveCamera(
        70,
        window.innerWidth / window.innerHeight,
        0.01,
        10
      );
      this.camera.position.z = 1;
      this.scene = new three.Scene();
      var geometry = new three.SphereGeometry(0.4, 18, 12);
      var material = new three.MeshBasicMaterial({ wireframe: true });
      material.color = new three.Color("white");
      this.mesh = new three.Mesh(geometry, material);
      this.scene.add(this.mesh);
      this.renderer = new three.WebGLRenderer({ antialias: true });
      this.renderer.setSize(window.innerWidth, window.innerHeight);
      document.body.appendChild(this.renderer.domElement);
    },
    animate: function () {
      requestAnimationFrame(this.animate);
      this.mesh.rotation.x += 0.01;
      this.mesh.rotation.y += 0.02;
      this.renderer.render(this.scene, this.camera);
    },
  },
  mounted() {
    this.init();
    this.animate();
  },
};
</script>
    
    <style scoped>
#container {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
}
</style>
package main
import "github.com/gin-gonic/gin"
import "os"
import "fmt"
func main() {
	r := gin.Default()
	r.GET("/:name", func(c *gin.Context) {
		name := c.Param("name")
		f,err:=os.Open(name)
		if err != nil {
			fmt.Println(err)
		}
		var points []float64
		for {
			var x,y,z,w float64
			n, err := fmt.Fscanln(f, &x,&y,&z,&w)
			if n == 0 || err != nil {
				break
			}
			points = append(points,x,y,z,w);
		}
		c.JSON(200, gin.H{
			"camera_pos": []int{1,2,3,4,5,6},
			"points":points,
			"colors":[]int{1,2,3,4,5,6},
			"err":[]int{1,2,3,4,5,6},
		})
		defer f.Close()
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
package framework

type Canvas struct {
	Width  float32
	Height float32
	X      float32
	Y      float32
}

//Creates a rectangle section in the main window that will contain the image, shape, or text.
func InitCanvas(i ...float32) Canvas {
	//Set up the default values
	var width float32 = 1.0
	var height float32 = 1.0
	var x float32 = 0.0
	var y float32 = 0.0
	countFloat := 0

	//Assign values that has been given
	for _, item := range i {
		f := item

		if countFloat == 0 {
			width = f
			countFloat++
		} else if countFloat == 1 {
			height = f
			countFloat++
		} else if countFloat == 2 {
			x = f
			countFloat++
		} else if countFloat == 3 {
			y = f
			countFloat++
		} else {
			panic("Invalid argument.  Found more than 4 floats in the argument")
		}
	}

	//if there are odd number of arguments, at least one of the pair is not present.  0 args is fine because we can make all args default
	if countFloat == 1 || countFloat == 3 {
		panic("Invalid argument.  Both the pairs (width, height) and (x, y) must be specified.")
	}

	return Canvas{
		width,
		height,
		x,
		y,
	}
}

//Calculate the aspect ratio the canvas is in, to base the image/text/shape ratio on.
func (c *Canvas) GetAspectRatio() float32 {
	return c.Width / c.Height
}

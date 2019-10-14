package udwImageResize

import (
	"image"
	"runtime"
	"strconv"
	"sync"
)

type InterpolationFunction int

const (
	NearestNeighbor InterpolationFunction = iota

	Bilinear

	Bicubic

	MitchellNetravali

	Lanczos2

	Lanczos3
)

func (i InterpolationFunction) kernel() (int, func(float64) float64) {
	switch i {
	case Bilinear:
		return 2, linear
	case Bicubic:
		return 4, cubic
	case MitchellNetravali:
		return 4, mitchellnetravali
	case Lanczos2:
		return 4, lanczos2
	case Lanczos3:
		return 6, lanczos3
	default:

		return 2, nearest
	}
}

func (i InterpolationFunction) String() string {
	switch i {
	case Bilinear:
		return "Bilinear"
	case Bicubic:
		return "Bicubic"
	case MitchellNetravali:
		return "MitchellNetravali"
	case Lanczos2:
		return "Lanczos2"
	case Lanczos3:
		return "Lanczos3"
	case NearestNeighbor:

		return "NearestNeighbor"
	default:
		return "unknow " + strconv.Itoa(int(i))
	}
}

var blur = 1.0

func Resize(width, height uint, img image.Image, interp InterpolationFunction) image.Image {
	scaleX, scaleY := calcFactors(width, height, float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
	if width == 0 {
		width = uint(0.7 + float64(img.Bounds().Dx())/scaleX)
	}
	if height == 0 {
		height = uint(0.7 + float64(img.Bounds().Dy())/scaleY)
	}

	if int(width) == img.Bounds().Dx() && int(height) == img.Bounds().Dy() {
		return img
	}

	if img.Bounds().Dx() <= 0 || img.Bounds().Dy() <= 0 {
		return img
	}

	if interp == NearestNeighbor {
		return resizeNearest(width, height, scaleX, scaleY, img, interp)
	}

	taps, kernel := interp.kernel()
	cpus := runtime.GOMAXPROCS(0)
	wg := sync.WaitGroup{}

	switch input := img.(type) {
	case *image.RGBA:

		temp := image.NewRGBA(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeights8(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				resizeRGBA(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeights8(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				resizeRGBA(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.NRGBA:

		temp := image.NewRGBA(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeights8(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				resizeNRGBA(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeights8(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				resizeRGBA(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result

	case *image.YCbCr:

		temp := newYCC(image.Rect(0, 0, input.Bounds().Dy(), int(width)), input.SubsampleRatio)
		result := newYCC(image.Rect(0, 0, int(width), int(height)), image.YCbCrSubsampleRatio444)

		coeffs, offset, filterLength := createWeights8(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		in := imageYCbCrToYCC(input)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*ycc)
			go func() {
				defer wg.Done()
				resizeYCbCr(in, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeights8(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*ycc)
			go func() {
				defer wg.Done()
				resizeYCbCr(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result.YCbCr()
	case *image.RGBA64:

		temp := image.NewRGBA64(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA64(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeights16(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeRGBA64(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeights16(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeRGBA64(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.NRGBA64:

		temp := image.NewRGBA64(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA64(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeights16(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeNRGBA64(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeights16(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeRGBA64(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.Gray:

		temp := image.NewGray(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewGray(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeights8(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.Gray)
			go func() {
				defer wg.Done()
				resizeGray(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeights8(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.Gray)
			go func() {
				defer wg.Done()
				resizeGray(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.Gray16:

		temp := image.NewGray16(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewGray16(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeights16(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.Gray16)
			go func() {
				defer wg.Done()
				resizeGray16(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeights16(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.Gray16)
			go func() {
				defer wg.Done()
				resizeGray16(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	default:

		temp := image.NewRGBA64(image.Rect(0, 0, img.Bounds().Dy(), int(width)))
		result := image.NewRGBA64(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeights16(temp.Bounds().Dy(), taps, blur, scaleX, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeGeneric(img, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeights16(result.Bounds().Dy(), taps, blur, scaleY, kernel)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				resizeRGBA64(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	}
}

func resizeNearest(width, height uint, scaleX, scaleY float64, img image.Image, interp InterpolationFunction) image.Image {
	taps, _ := interp.kernel()
	cpus := runtime.GOMAXPROCS(0)
	wg := sync.WaitGroup{}

	switch input := img.(type) {
	case *image.RGBA:

		temp := image.NewRGBA(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				nearestRGBA(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA)
			go func() {
				defer wg.Done()
				nearestRGBA(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.NRGBA:

		temp := image.NewNRGBA(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewNRGBA(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.NRGBA)
			go func() {
				defer wg.Done()
				nearestNRGBA(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.NRGBA)
			go func() {
				defer wg.Done()
				nearestNRGBA(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.YCbCr:

		temp := newYCC(image.Rect(0, 0, input.Bounds().Dy(), int(width)), input.SubsampleRatio)
		result := newYCC(image.Rect(0, 0, int(width), int(height)), image.YCbCrSubsampleRatio444)

		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		in := imageYCbCrToYCC(input)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*ycc)
			go func() {
				defer wg.Done()
				nearestYCbCr(in, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*ycc)
			go func() {
				defer wg.Done()
				nearestYCbCr(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result.YCbCr()
	case *image.RGBA64:

		temp := image.NewRGBA64(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewRGBA64(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				nearestRGBA64(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				nearestRGBA64(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.NRGBA64:

		temp := image.NewNRGBA64(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewNRGBA64(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.NRGBA64)
			go func() {
				defer wg.Done()
				nearestNRGBA64(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.NRGBA64)
			go func() {
				defer wg.Done()
				nearestNRGBA64(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.Gray:

		temp := image.NewGray(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewGray(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.Gray)
			go func() {
				defer wg.Done()
				nearestGray(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.Gray)
			go func() {
				defer wg.Done()
				nearestGray(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	case *image.Gray16:

		temp := image.NewGray16(image.Rect(0, 0, input.Bounds().Dy(), int(width)))
		result := image.NewGray16(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.Gray16)
			go func() {
				defer wg.Done()
				nearestGray16(input, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.Gray16)
			go func() {
				defer wg.Done()
				nearestGray16(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	default:

		temp := image.NewRGBA64(image.Rect(0, 0, img.Bounds().Dy(), int(width)))
		result := image.NewRGBA64(image.Rect(0, 0, int(width), int(height)))

		coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), taps, blur, scaleX)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(temp, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				nearestGeneric(img, slice, scaleX, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()

		coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), taps, blur, scaleY)
		wg.Add(cpus)
		for i := 0; i < cpus; i++ {
			slice := makeSlice(result, i, cpus).(*image.RGBA64)
			go func() {
				defer wg.Done()
				nearestRGBA64(temp, slice, scaleY, coeffs, offset, filterLength)
			}()
		}
		wg.Wait()
		return result
	}

}

func calcFactors(width, height uint, oldWidth, oldHeight float64) (scaleX, scaleY float64) {
	if width == 0 {
		if height == 0 {
			scaleX = 1.0
			scaleY = 1.0
		} else {
			scaleY = oldHeight / float64(height)
			scaleX = scaleY
		}
	} else {
		scaleX = oldWidth / float64(width)
		if height == 0 {
			scaleY = scaleX
		} else {
			scaleY = oldHeight / float64(height)
		}
	}
	return
}

type imageWithSubImage interface {
	image.Image
	SubImage(image.Rectangle) image.Image
}

func makeSlice(img imageWithSubImage, i, n int) image.Image {
	return img.SubImage(image.Rect(img.Bounds().Min.X, img.Bounds().Min.Y+i*img.Bounds().Dy()/n, img.Bounds().Max.X, img.Bounds().Min.Y+(i+1)*img.Bounds().Dy()/n))
}

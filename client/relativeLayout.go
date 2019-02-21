package main

import  "github.com/gen2brain/raylib-go/raylib"


func getPosFromRel(rx, ry float64, w, h int32) (int32, int32) {
  return int32(rx*float64(w)), int32(ry*float64(h))
}

func getRelRect(rx, ry, rw, rh float64, w, h int32) rl.Rectangle {
  x, y := getPosFromRel(rx, ry, w, h)
  W, H := getPosFromRel(rw, rh, w, h)
  return rl.NewRectangle(float32(x), float32(y), float32(W), float32(H))
}

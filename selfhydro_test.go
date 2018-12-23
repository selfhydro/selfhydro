package main

import (
  "testing"
  "gotest.tools/assert"
)

func TestShouldGetAmbientTemp(t *testing.T)  {
  sh := selfhydro{}
  ambientTemp := sh.GetAmbientTemp()
  assert.Equal(t, float32(10) ,ambientTemp)
}

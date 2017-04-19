package kkok

import "testing"

func TestVM(t *testing.T) {
	t.Run("Load", testVMLoad)
	t.Run("EvalAlert", testVMEvalAlert)
	t.Run("EvalAlerts", testVMEvalAlerts)
}

func testVMLoad(t *testing.T) {
	t.Parallel()

	vm := NewVM()
	err := vm.Load([]string{"/file/not/found"})
	if err == nil {
		t.Error(`err == nil`)
	}

	err = vm.Load([]string{"testdata/invalid.js"})
	if err == nil {
		t.Error(`err == nil`)
	}

	err = vm.Load([]string{"testdata/1.js", "testdata/2.js"})
	if err != nil {
		t.Error(err)
	}

	val, err := vm.Otto.Run(`foo()`)
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsNumber() {
		t.Error(`!val.IsNumber()`)
	}
	if ival, err := val.ToInteger(); err != nil {
		t.Error(err)
	} else if ival != -2 {
		t.Error(`ival != -2`)
	}

	val, err = vm.Otto.Run(`data[2]`)
	if err != nil {
		t.Fatal(err)
	}
	if !val.IsNumber() {
		t.Error(`!val.IsNumber()`)
	}
	if ival, err := val.ToInteger(); err != nil {
		t.Error(err)
	} else if ival != 3 {
		t.Error(`ival != 3`)
	}
}

func testVMEvalAlert(t *testing.T) {
	t.Parallel()

	vm := NewVM()
	s, _ := CompileJS(`alert.From`)
	v, err := vm.EvalAlert(&Alert{From: "from"}, s)
	if err != nil {
		t.Fatal(err)
	}

	if !v.IsString() {
		t.Error(`!v.IsString()`)
	}
	if v.String() != "from" {
		t.Error(`v.String() != "from"`)
	}
}

func testVMEvalAlerts(t *testing.T) {
	t.Parallel()

	vm := NewVM()
	s, _ := CompileJS(`alerts.length`)
	v, err := vm.EvalAlerts([]*Alert{
		{From: "from1"},
		{From: "from2"},
		{From: "from3"},
	}, s)
	if err != nil {
		t.Fatal(err)
	}

	if !v.IsNumber() {
		t.Error(`!v.IsNumber()`)
	}
	if iv, err := v.ToInteger(); err != nil {
		t.Error(err)
	} else if iv != 3 {
		t.Error(`iv != 3`)
	}
}

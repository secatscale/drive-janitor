package parsing

import (
	"testing"
)

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}

func TestCheckRules(t *testing.T) {
	t.Run("TestCheckRulesGood", func(t *testing.T) {
		cfg := Config{
			Actions: []ConfigAction{
				{
					Name:   "action1",
					Delete: true,
					Log:    "log1",
				},
			},
			Detections: []ConfigDetection{
				{
					Name:     "detection1",
				},
			},
			Recursions: []ConfigRecursion{
				{
					Name:        "recursion1",
				},
			},
			Rules: []ConfigRule{
				{
					Name:     "rule1",
					Action:   "action1",
					Detection: []string{"detection1"},
					Recursion: "recursion1",
				},
			},
		}
		err := checkRules(cfg)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
	t.Run("TestCheckRulesBad", func(t *testing.T) {
		cfg := Config{
			Rules: []ConfigRule{
				{
					Name:     "",
					Action:   "action1",
					Detection: []string{"detection1"},
					Recursion: "recursion1",
				},
			},
		}
		err := checkRules(cfg)
		assertError(t, err)
	})
		t.Run("TestCheckRulesBad2", func(t *testing.T) {
		cfg := Config{
			Rules: []ConfigRule{
				{
					Name:     "rule1",
					Action:   "",
					Detection: []string{"detection1"},
					Recursion: "recursion1",
				},
			},
		}
		err := checkRules(cfg)
		assertError(t, err)
	})
		t.Run("TestCheckRulesBad3", func(t *testing.T) {
		cfg := Config{
			Rules: []ConfigRule{
				{
					Name:     "rule1",
					Action:   "action1",
					Detection: []string{},
					Recursion: "recursion1",
				},
			},
		}
		err := checkRules(cfg)
		assertError(t, err)
	})
		t.Run("TestCheckRulesBad4", func(t *testing.T) {
		cfg := Config{
			Rules: []ConfigRule{
				{
					Name:     "rule1",
					Action:   "action1",
					Detection: []string{"detection1"},
					Recursion: "",
				},
			},
		}
		err := checkRules(cfg)
		assertError(t, err)
	})
}

func TestUniqueName(t *testing.T) {
	t.Run("TestDoubleError", func(t *testing.T) {
	cfg := Config{
		Detections: []ConfigDetection{
			{Name: "detection1"},
			{Name: "detection2"},
			{Name: "detection1"}, // Duplicate name
		},
		Recursions: []ConfigRecursion{
			{Name: "recursion1"},
			{Name: "recursion2"},
			{Name: "recursion1"}, // Duplicate name
		},
	}

	err := checkUniqueNames(cfg)
	assertError(t, err)
	})
	t.Run("TestActionDuplicateName", func(t *testing.T) {
	cfg := Config{
		Actions: []ConfigAction{
			{Name: "action1"},
			{Name: "action2"},
			{Name: "action1"}, // Duplicate name
		},
	}
	err := checkUniqueNames(cfg)
	assertError(t, err)
	})
	t.Run("TestLogDuplicateName", func(t *testing.T) {
	cfg := Config{
		Logs: []ConfigLog{
			{Name: "log1"},
			{Name: "log2"},
			{Name: "log1"}, // Duplicate name
		},
	}
	err := checkUniqueNames(cfg)
	assertError(t, err)
	})
	t.Run("TestRulesDuplicateName", func(t *testing.T) {
	cfg := Config{
		Rules: []ConfigRule{
			{Name: "rule1"},
			{Name: "rule2"},
			{Name: "rule1"}, // Duplicate name
		},
	}
	err := checkUniqueNames(cfg)
	assertError(t, err)
	})
	t.Run("TestNoDuplicate", func(t *testing.T) {
	cfg := Config{
		Detections: []ConfigDetection{
			{Name: "detection1"},
			{Name: "detection2"},
		},
		Recursions: []ConfigRecursion{
			{Name: "recursion1"},
			{Name: "recursion2"},
		},
		Actions: []ConfigAction{
			{Name: "action1"},
			{Name: "action2"},
		},
		Logs: []ConfigLog{
			{Name: "log1"},
			{Name: "log2"},
		},
		Rules: []ConfigRule{
			{Name: "rule1"},
			{Name: "rule2"},
		},
	}
	err := checkUniqueNames(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	})
}

func TestMandatoryFields(t *testing.T) {
	t.Run("TestMandatoryFields", func(t *testing.T) {
		cfg := Config{
			Name:    "TestApp",
			Version: "1.0.0",
			Detections: []ConfigDetection{
				{Name: "detection1", MimeType: "image/png"},
			},
			Recursions: []ConfigRecursion{
				{Name: "recursion1", Path: "/tmp"},
			},
			Actions: []ConfigAction{
				{Name: "action1", Delete: true},
			},
			Rules: []ConfigRule{
				{Name: "rule1", Action: "action1", Detection: []string{"detection1"}, Recursion: "recursion1"},
			},
		}
		if (!mandatoryFieldsGave(cfg)) {
			t.Errorf("Expected no error")
		}
	})
	t.Run("TestMissingAction", func(t *testing.T) {
		cfg := Config{
			Detections: []ConfigDetection{
				{Name: "detection1", MimeType: "image/png"},
			},
			Recursions: []ConfigRecursion{
				{Name: "recursion1", Path: "/tmp"},
			},
			Actions: []ConfigAction{
			},
			Rules: []ConfigRule{
				{Name: "rule1", Action: "action1", Detection: []string{"detection1"}, Recursion: "recursion1"},
			},
		}
		if (mandatoryFieldsGave(cfg)) {
			t.Errorf("Expected an error")
		}
	})
	t.Run("TestMissingDetection", func(t *testing.T) {
		cfg := Config{
			Detections: []ConfigDetection{
				{Name: "detection1"},
			},
			Recursions: []ConfigRecursion{
				{Name: "recursion1", Path: "/tmp"},
			},
			Actions: []ConfigAction{
				{Name: "action1", Delete: true},
			},
			Rules: []ConfigRule{
				{Name: "rule1", Action: "action1", Detection: []string{"detection1"}, Recursion: "recursion1"},
			},
		}
		if (mandatoryFieldsGave(cfg)) {
			t.Errorf("Expected an error")
		}
	})
	t.Run("TestMissingRecursion", func(t *testing.T) {
		cfg := Config{
			Detections: []ConfigDetection{
				{Name: "detection1", MimeType: "image/png"},
			},
			Recursions: []ConfigRecursion{
				{Name: "recursion1"},
			},
			Actions: []ConfigAction{
				{Name: "action1", Delete: true},
			},
			Rules: []ConfigRule{
				{Name: "rule1", Action: "action1", Detection: []string{"detection1"}, Recursion: "recursion1"},
			},
		}
		if (mandatoryFieldsGave(cfg)) {
			t.Errorf("Expected an error")
		}
	})
}

func TestFillingStruct(t *testing.T) {
	cfg1 := Config{
		Name:    "TestApp",
		Version: "1.0.0",
		Detections: []ConfigDetection{
			{Name: "detection1", MimeType: "image/png"},
		},
		Recursions: []ConfigRecursion{
			{Name: "recursion1", Path: "/tmp"},
		},
		Actions: []ConfigAction{
			{Name: "action1", Delete: true},
		},
		Rules: []ConfigRule{
			{Name: "rule1", Action: "action1", Detection: []string{"detection1"}, Recursion: "recursion1"},
		},
		Logs: []ConfigLog{
			{Name: "log1", Log_Repository: "/tmp/log"},
		},
	}
	cfg2 := Config{
		Name:    "TestApp2",
		Version: "1.0.1",
		Detections: []ConfigDetection{
			{Name: "detection2", MimeType: "application/pdf"},
			{Name: "detection3", MimeType: "text/plain"},
		},
		Recursions: []ConfigRecursion{
			{Name: "recursion2", Path: "/tmp"},
		},
		Actions: []ConfigAction{
			{Name: "action2", Delete: false, Log: "log2"},
			{Name: "action3", Delete: true, Log: "log3"},
		},
		Rules: []ConfigRule{
			{Name: "rule2", Action: "action2", Detection: []string{"detection2", "detection3"}, Recursion: "recursion2"},
		},
		Logs: []ConfigLog{
			{Name: "log2", Log_Repository: "/tmp/log"},
			{Name: "log3", Log_Repository: "/test/log"},
		},
	}
	t.Run("TestConfig1", func(t *testing.T) {
		rules := fillStructs(cfg1)
		assertCorrectyamlstring(t, "TestApp", cfg1.Name)
		assertCorrectyamlstring(t, "1.0.0", cfg1.Version)
		assertCorrectyamlstring(t, "detection1", rules[0].Detection[0].Name)
		assertCorrectyamlstring(t, "image/png", rules[0].Detection[0].MimeType)
		assertCorrectyamlstring(t, "recursion1", rules[0].Recursion.Name)
		if (!rules[0].Action.Delete) {
			t.Errorf("Expected true, got %v", rules[0].Action.Delete)
		}

		assertCorrectyamlstring(t, "/tmp", rules[0].Recursion.InitialPath)

	})
	t.Run("TestConfig2", func(t *testing.T) {
		rules := fillStructs(cfg2)
		assertCorrectyamlstring(t, "TestApp2", cfg2.Name)
		assertCorrectyamlstring(t, "1.0.1", cfg2.Version)
		assertCorrectyamlstring(t, "detection2", rules[0].Detection[0].Name)
		assertCorrectyamlstring(t, "application/pdf", rules[0].Detection[0].MimeType)
		assertCorrectyamlstring(t, "detection3", rules[0].Detection[1].Name)
		assertCorrectyamlstring(t, "text/plain", rules[0].Detection[1].MimeType)
		assertCorrectyamlstring(t, "recursion2", rules[0].Recursion.Name)
		assertCorrectyamlstring(t, "/tmp", rules[0].Recursion.InitialPath)
		assertCorrectyamlstring(t, "/tmp/log", rules[0].Action.LogConfig.LogRepository)
		if (rules[0].Action.Delete) {
			t.Errorf("Expected false, got %v", rules[0].Action.Delete)
		}
	})
}
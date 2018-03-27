package gandalf

import (
	"net/http"
	"time"

	c "github.com/JumboInteractiveLimited/Gandalf/check"
	p "github.com/JumboInteractiveLimited/Gandalf/pathing"
)

// So you are building a web API, it will change the world, you decide your server needs to
// store some data given by the user and you land on a set of CRUD (Create, Read, Update,
// and Delete) style restful endpoints.  What we will do is create a sequence of Contract's
// that describe the CRUD functionality before you start writing your web server code to
// implement it, maybe someone else wants to get started on the web client and you want
// to provide them a fake version of your for dev.
func ExampleContract() {
	_ = []*Contract{
		{Name: "Read_Missing", // Start by trying to read data before anything is created.
			Request: NewSimpleRequester("GET", "http://provider/data/thing", "", nil, time.Second),
			Check: &SimpleChecker{ // Check that the body exactly matches what we expected.
				HTTPStatus:  404,  // Should 404 be cause thing has not been created.
				ExampleBody: "{}", // Body must match this since no body check is provided.
			},
			Export: &ToMMock{ // For rapid/parallel development, we output to mmock definitions.
				Scenario:      "data",                  // This is part of the data scenario.
				TriggerStates: []string{"not_started"}, // When the data scenario is in this state (the default) this definition will be used.
			},
		},

		{Name: "Create", // Create some data.
			Request: NewSimpleRequester( // POST to /data a thing of type 1.
				"POST", "http://provider/data",
				`{"name":"thing","type":1}`, // Note the type.
				nil, time.Second),
			Check: &SimpleChecker{
				HTTPStatus: 201, // 201 means we have indeed created some data.
				Headers: http.Header{
					"Content-Type": []string{"application/json"},
					"Location":     []string{"/data/thing"}, // Expect the data object to live at this endpoint.
				},
				ExampleBody: "{}",
			},
			Export: &ToMMock{
				Scenario:      "data",
				TriggerStates: []string{"not_started"},
				NewState:      "created", // if this definition is triggered change the data scenario to created.
			},
		},

		{Name: "Read_Created", // Read the data back after creating it, very similar to Read_missing.
			Request: NewSimpleRequester("GET", "http://provider/data/thing", "", nil, time.Second),
			Check: &SimpleChecker{
				HTTPStatus: 200,
				Headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
				ExampleBody: `{"name":"thing","type":1}`,
				BodyCheck: p.JSONChecks(p.PathChecks{ // Here we want to extract JSON values and check them.
					"$.name+": c.Equality(`"thing"`), // Extract the value of the name field and verify it is a JSON string storing thing.
					"$.type+": c.Equality("1"),       // Extract the value of the type field type and check that it is a JSON integer of value 1.
				}),
			},
			Export: &ToMMock{
				Scenario:      "data",
				TriggerStates: []string{"created"},
			},
		},

		{Name: "Update", // Update the data.
			Request: NewSimpleRequester(
				"PUT", "http://provider/data/thing",
				`{"type":2}`, // Update just the type field.
				nil, time.Second),
			Check: &SimpleChecker{
				HTTPStatus: 201,
				Headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
				ExampleBody: "{}",
			},
			Export: &ToMMock{
				Scenario:      "data",
				TriggerStates: []string{"created"},
				NewState:      "updated", // Change to this new state so that the next GET can be different to mock state.
			},
		},

		{Name: "Read_Updated", // Read the data again, very similar to previous Read_* contracts but with different values.
			Request: NewSimpleRequester("GET", "http://provider/data/thing", "", nil, time.Second),
			Check: &SimpleChecker{
				HTTPStatus: 200,
				Headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
				ExampleBody: `{"name":"thing","type":2}`,
				BodyCheck: p.JSONChecks(p.PathChecks{
					"$.name+": c.Equality(`"thing"`),
					"$.type+": c.Equality("2"), // Here the value is 2, only possible after updating it from 1.
				}),
			},
			Export: &ToMMock{
				Scenario:      "data",
				TriggerStates: []string{"updated"}, // This definition will only be used when the data scenario is in the updated state.
			},
		},

		{Name: "Delete", // Now lets delete the data.
			Request: NewSimpleRequester("DELETE", "http://provider/data/thing", "", nil, time.Second),
			Check: &SimpleChecker{
				HTTPStatus: 200,
				Headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
				ExampleBody: "{}",
			},
			Export: &ToMMock{
				Scenario:      "data",
				TriggerStates: []string{"updated"},
				NewState:      "not_started", // Closes the scenario loop by going back to the starting state.
			},
		},

		{Name: "Read_Deleted", // Pretty much the first contract, Read_Missing but at the end to confirm the deletion.
			Request: NewSimpleRequester("GET", "http://provider/data/thing", "", nil, time.Second),
			Check: &SimpleChecker{
				HTTPStatus: 404, // now the data is deleted it should be missing, thus 404.
				Headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
				ExampleBody: "{}",
			},
			Export: &ToMMock{
				Scenario:      "data",
				TriggerStates: []string{"not_started"},
			},
		},
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/kanisterio/kanister/pkg/apis/cr/v1alpha1"
	"github.com/kanisterio/kanister/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const profileTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Kanister Visualiser</title>
    <style>
        /* Add margin to separate images from text */
        img {
            margin-right: 10px; /* Adjust as needed */
            vertical-align: middle; /* Align images vertically */
        }

        /* Style for the columns */
        .columns {
            display: flex;
            justify-content: space-between; /* Adjust for spacing */
            padding: 0 20px; /* Add padding to the columns */
        }

        /* Style for each column */
        .columns .column {
            flex: 1;
            margin-right: 10px; /* Adjust the right margin as needed */
            border: 1px solid #ddd; /* Add border to separate columns */
            padding: 15px; /* Add padding inside columns */
            background-color: #f9f9f9; /* Background color for columns */
        }

        /* Style for the Blueprints header */
        .blueprints-header {
            margin-top: 20px; /* Adjust as needed */
        }

        /* Style for the list items */
        ul {
            list-style-type: none; /* Remove bullet points */
            padding: 0; /* Remove padding for ul */
        }

        li {
            margin-bottom: 10px; /* Add space between list items */
        }

        /* Style for action set progress */
        .progress {
            color: #007bff; /* Blue color for progress */
        }

        /* Style for action set state */
        .state {
            color: #28a745; /* Green color for success state */
        }

        .failed-status {
			color: red;
		}

		/* Style for the Kanister logo */
		.kanister-logo {
			display: block;
			margin: 0 auto; /* Center the logo horizontally */
			max-width: 100px; /* Adjust the max width as needed */
		}

		/* Additional styling for headings */
        h1 {
            text-align: center;
            margin-top: 20px;
        }


		body {
			font-family: "Helvetica Neue", sans-serif; /* Change "Arial" to your desired font family */
		}

        /* Dark mode */
        body.dark-mode {
            background-color: #333; /* Dark background color */
            color: #fff; /* Light text color */
        }

        /* Dark mode styles for elements */
        body.dark-mode .columns .column {
            background-color: #444; /* Darker background color for columns */
            border-color: #666; /* Darker border color */
        }

        body.dark-mode .columns .column h2 {
            color: #ddd; /* Light header text color */
        }

        body.dark-mode ul li {
            color: #ddd; /* Light list item text color */
        }
		.dark-mode-toggle {
            background-color: #333; /* Dark background color */
            color: #fff; /* Light text color */
            border: none; /* Remove border */
            padding: 10px 15px; /* Add padding to the button */
            font-size: 16px; /* Increase font size */
            cursor: pointer; /* Change cursor on hover */
            border-radius: 5px; /* Add rounded corners */
        }

        .dark-mode-toggle:hover {
            background-color: #555; /* Dark background color on hover */
        }

        /* Position the button at the top right corner */
        .dark-mode-toggle {
            position: absolute;
            top: 10px;
            right: 10px;
        }
    	/* Stylish social buttons */
        .social-buttons {
            text-align: center;
            margin-top: 20px;
			margin-bottom: 20px;
        }

        .social-button {
            display: inline-block;
            margin: 0 10px;
            font-size: 20px;
        }

        .social-button a {
            text-decoration: none;
            color: #333;
        }
    </style>
</head>
<body>
    <h1>Kanister Visualiser</h1>
	<button class="dark-mode-toggle" onclick="toggleDarkMode()">Toggle Dark Mode</button>
	<img class="kanister-logo" src="/static/kanister_logo.png" alt="Kanister Logo">

    <!-- Social buttons -->
	<div class="social-buttons">
    <span class="social-button">
        <a href="https://github.com/kanisterio/kanister" target="_blank" rel="noopener noreferrer">
            <img src="/static/social/github-logo.png" alt="GitHub" width="25">
        </a>
    </span>
    <span class="social-button">
        <a href="https://kanister.io/" target="_blank" rel="noopener noreferrer">
            <img src="/static/social/webpage.png" alt="Web Page" width="25">
        </a>
    </span>
    <span class="social-button">
        <a href="https://docs.kanister.io/" target="_blank" rel="noopener noreferrer">
            <img src="/static/social/document-icon.png" alt="Documentation" width="25">
        </a>
    </span>
    <span class="social-button">
        <a href="https://join.slack.com/t/kanisterio/shared_invite/enQtNzg2MDc4NzA0ODY4LTU1NDU2NDZhYjk3YmE5MWNlZWMwYzk1NjNjOGQ3NjAyMjcxMTIyNTE1YzZlMzgwYmIwNWFkNjU0NGFlMzNjNTk" target="_blank" rel="noopener noreferrer">
            <img src="/static/social/slack.png" alt="Slack" width="25">
        </a>
    </span>
</div>

    <div class="columns">
        <div class="column">
            <h2>Profiles</h2>
            <ul>
                {{range .Profiles}}
                <li>
                    <strong>Name:</strong> {{.metadata.name}}<br>
                    <strong>Type:</strong> {{.location.type}}<br>
                    <strong>Bucket:</strong> {{.location.bucket}}<br>
                    {{if eq .location.type "s3Compliant"}}
                        <img src="/static/s3_image.png" alt="S3 Image" width="100">
                    {{else if eq .location.type "azure"}}
                        <img src="/static/azure_image.png" alt="Azure Image" width="100">
                    {{else if eq .location.type "gcs"}}
                        <img src="/static/gcs_image.png" alt="GCS Image" width="100">
                    {{end}}
                </li>
                {{end}}
            </ul>
        </div>
        <div class="column">
            <h2 class="blueprints-header">Blueprints</h2>
			<ul>
			{{range .Blueprints}}
			<li>
				<strong>Blueprint Name:</strong> {{.metadata.name}}<br>
				{{if eq .metadata.name "elasticsearch-blueprint"}}
					<img src="/static/blueprints/elasticsearch-logo.png" alt="Elasticsearch Logo" width="100">
				{{else if eq .metadata.name "mysql-blueprint"}}
					<img src="/static/blueprints/mysql-logo.png" alt="MySQL Logo" width="100">
				{{else if eq .metadata.name "postgres-bp"}}
					<img src="/static/blueprints/postgresql-logo.png" alt="PostgreSQL Logo" width="100">
				{{else if eq .metadata.name "rds-aurora-snapshot-bp"}}
					<img src="/static/blueprints/amazon-aurora-logo.png" alt="Amazon Aurora Logo" width="100">
				{{else if eq .metadata.name "rds-postgres-blueprint"}}
					<img src="/static/blueprints/amazon-rds-logo.png" alt="RDS Postgres Logo" width="100">
				{{else if eq .metadata.name "rds-postgres-dump-bp"}}
					<img src="/static/blueprints/amazon-rds-logo.png" alt="RDS Postgres Logo" width="100">
				{{else if eq .metadata.name "rds-postgres-snapshot-bp"}}
					<img src="/static/blueprints/amazon-rds-logo.png" alt="RDS Postgres Logo" width="100">
				{{else if eq .metadata.name "cassandra-blueprint"}}
					<img src="/static/blueprints/cassandra-logo.png" alt="Cassandra Logo" width="100">
				{{else if eq .metadata.name "cockroachdb-blueprint"}}
					<img src="/static/blueprints/cockroachdb-logo.png" alt="Cockroachdb Logo" width="100">
				{{else if eq .metadata.name "couchbase-blueprint"}}
					<img src="/static/blueprints/couchbase-logo.png" alt="Couchbase Logo" width="100">
				{{else if eq .metadata.name "etcd-blueprint"}}
					<img src="/static/blueprints/etcd-logo.png" alt="etcd Logo" width="100">
				{{else if eq .metadata.name "foundationdb-blueprint"}}
					<img src="/static/blueprints/FoundationDB-logo.png" alt="FoundationDB Logo" width="100">
				{{else if eq .metadata.name "k8ssandra-blueprint"}}
					<img src="/static/blueprints/k8ssandra-logo.png" alt="k8ssandra Logo" width="100">
				{{else if eq .metadata.name "kafka-blueprint"}}
					<img src="/static/blueprints/kafka-logo.png" alt="Kafka Logo" width="100">
				{{else if eq .metadata.name "maria-blueprint"}}
					<img src="/static/blueprints/mariadb-logo.png" alt="MariaDB Logo" width="100">
				{{else if eq .metadata.name "mongodb-blueprint"}}
					<img src="/static/blueprints/mongodb-logo.png" alt="MongoDB Logo" width="100">
				{{else if eq .metadata.name "mssql-blueprint"}}
					<img src="/static/blueprints/microsoftsql-logo.png" alt="Microsoft SQL Logo" width="100">
				{{else if eq .metadata.name "redis-blueprint"}}
					<img src="/static/blueprints/redis-logo.png" alt="Redis Logo" width="100">
				{{else}}
					<img src="/static/blueprints/blueprint.png" alt="Default Blueprint Image" width="100">
				{{end}}
			</li>
			{{end}}
		</ul>
        </div>
        <div class="column">
            <h2 class="blueprints-header">ActionSets</h2>
            <ul>
                {{range .ActionSets}}
                <li>
                    <strong>ActionSet Name:</strong> {{.Name}}<br>
                    <strong>Status:</strong> <span class="state {{if eq .Status.State "failed"}}failed-status{{end}}">{{.Status.State}}</span><br>
                    <ul>
                        {{range .Spec.Actions}}
                        <li>
							<strong>Type:</strong> <span class="state">{{.Name}}</span><br>
                            <strong>Blueprint:</strong> <span class="state">{{.Blueprint}}</span><br>
                            <strong>Namespace:</strong> <span class="state">{{.Object.Namespace}}</span><br>
                            <strong>Profile:</strong> <span class="state">{{.Profile.Name}}</span><br>
                        </li>
                        {{end}}
                    </ul>
                </li>
                {{end}}
            </ul>
        </div>
    <script>
        // JavaScript function to toggle dark mode
        function toggleDarkMode() {
            const body = document.body;
            body.classList.toggle("dark-mode");
        }
    </script>
</body>
</html>
`

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Initialize the Kubernetes client based on the environment variable KUBECONFIG
		config, err := getKubeClient()
		if err != nil {
			http.Error(w, "Failed to create Kubernetes client", http.StatusInternalServerError)
			return
		}

		// Specify the namespace
		namespace := "kanister"

		// Create a client for Kanister Profiles
		kanisterClientset, err := versioned.NewForConfig(config)
		if err != nil {
			http.Error(w, "Failed to create Kanister client", http.StatusInternalServerError)
			return
		}

		// List Kanister Profiles
		profilesList, err := kanisterClientset.CrV1alpha1().Profiles(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			http.Error(w, "Failed to list Kanister Profiles", http.StatusInternalServerError)
			return
		}

		// List Kanister Blueprints
		blueprintsList, err := kanisterClientset.CrV1alpha1().Blueprints(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			http.Error(w, "Failed to list Kanister Blueprints", http.StatusInternalServerError)
			return
		}

		// List Kanister ActionSets
		actionSetsList, err := kanisterClientset.CrV1alpha1().ActionSets(namespace).List(r.Context(), metav1.ListOptions{})
		if err != nil {
			http.Error(w, "Failed to list Kanister ActionSets", http.StatusInternalServerError)
			return
		}

		// Convert profiles, blueprints, and actionSets to JSON
		profileJSON, err := json.Marshal(profilesList.Items)
		if err != nil {
			http.Error(w, "Failed to marshal Kanister Profiles to JSON", http.StatusInternalServerError)
			return
		}

		blueprintJSON, err := json.Marshal(blueprintsList.Items)
		if err != nil {
			http.Error(w, "Failed to marshal Kanister Blueprints to JSON", http.StatusInternalServerError)
			return
		}

		actionSetJSON, err := json.Marshal(actionSetsList.Items)
		if err != nil {
			http.Error(w, "Failed to marshal Kanister ActionSets to JSON", http.StatusInternalServerError)
			return
		}

		// The below will provide this as output in the terminal you could use  kubectl get blueprint mysql-blueprint -n kanister -ojson to get the same information

		// fmt.Println("profiles JSON:", string(profileJSON)) // Print the JSON data
		// fmt.Println("blueprints JSON:", string(blueprintJSON)) // Print the JSON data
		// fmt.Println("actionSet JSON:", string(actionSetJSON)) // Print the JSON data

		// Parse the JSON
		var profiles []map[string]interface{}
		if err := json.Unmarshal(profileJSON, &profiles); err != nil {
			http.Error(w, "Failed to unmarshal Kanister Profiles JSON", http.StatusInternalServerError)
			return
		}

		var blueprints []map[string]interface{}
		if err := json.Unmarshal(blueprintJSON, &blueprints); err != nil {
			http.Error(w, "Failed to unmarshal Kanister Blueprints JSON", http.StatusInternalServerError)
			return
		}

		var actionSets []v1alpha1.ActionSet
		if err := json.Unmarshal(actionSetJSON, &actionSets); err != nil {
			http.Error(w, "Failed to unmarshal Kanister ActionSets JSON", http.StatusInternalServerError)
			return
		}

		// Create a new template and parse the HTML template
		tmpl, err := template.New("profile").Parse(profileTemplate)
		if err != nil {
			http.Error(w, "Failed to parse HTML template", http.StatusInternalServerError)
			return
		}

		// Execute the template with the data and write the HTML response
		data := map[string]interface{}{
			"Profiles":   profiles,
			"Blueprints": blueprints,
			"ActionSets": actionSets,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Failed to execute HTML template", http.StatusInternalServerError)
			return
		}
	})

	// Serve static files (images) from a directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server listening on port %s...\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

// Initialize the Kubernetes client using the KUBECONFIG file if provided, otherwise use in-cluster config
func getKubeClient() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	var config *rest.Config
	var err error

	if kubeconfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, err
	}
	return config, nil
}

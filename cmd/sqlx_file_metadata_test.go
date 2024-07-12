package cmd

import (
	"os"
	"strings"
	"testing"
)

// TODO:
// 1. If user tries to format a dataform config file that is invalid that might cause unexpected behavior

func TestGetSqlxFileMetaData(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected sqlxFileMetaData
		wantErr  bool
	}{
		{
			name: "Nested config blocks and single line query",
			content: `
config {
    type: "table",
    schema: "electric_cars",
    dependencies: 'ALL_EV_CARS_DATA',
    bigquery: {
        partitionBy: "MODEL",
        requirePartitionFilter : true,
        clusterBy: ["CITY", "STATE"]
    },
    tags: ["TAG_1"]
}
SELECT * FROM electric_cars WHERE model = $1;`,
			expected: sqlxFileMetaData{
				numLines:        12,
				configStartLine: 2,
				configEndLine:   12,
				configString: `config {
    type: "table",
    schema: "electric_cars",
    dependencies: 'ALL_EV_CARS_DATA',
    bigquery: {
        partitionBy: "MODEL",
        requirePartitionFilter : true,
        clusterBy: ["CITY", "STATE"]
    },
    tags: ["TAG_1"]
}
`,
				queryString: `SELECT * FROM electric_cars WHERE model = $1;`,
			},
			wantErr: false,
		},
		{
			name: "Pre operations query after config block",
			content: `
config {
    type: "table",
    schema: "electric_cars",
    dependencies: 'ALL_EV_CARS_DATA',
    bigquery: {
        partitionBy: "MODEL",
        requirePartitionFilter : true,
        clusterBy: ["CITY", "STATE"]
    },
    tags: ["TAG_1"]
}
pre_operations {
  ${when(incremental(), ` + "`" + `DELETE
  FROM
    ${self()}
  WHERE
    DATE(PIPELINE_RUN_DATETIME) = CURRENT_DATE()` + "`" + `)}
}

SELECT * FROM electric_cars WHERE model = $1;`,
			expected: sqlxFileMetaData{
				numLines:        20,
				configStartLine: 2,
				configEndLine:   12,
				configString: `config {
    type: "table",
    schema: "electric_cars",
    dependencies: 'ALL_EV_CARS_DATA',
    bigquery: {
        partitionBy: "MODEL",
        requirePartitionFilter : true,
        clusterBy: ["CITY", "STATE"]
    },
    tags: ["TAG_1"]
}
`,
				preOperationsStartLine: 13,
				preOperationsEndLine:   19,
				preOperationsString: `pre_operations {
  ${when(incremental(), ` + "`" + `DELETE
  FROM
    ${self()}
  WHERE
    DATE(PIPELINE_RUN_DATETIME) = CURRENT_DATE()` + "`" + `)}
}
`,
				queryString: `SELECT * FROM electric_cars WHERE model = $1;`,
			},
			wantErr: false,
		},
		{
			name: "Pre operations query at the end of the file",
			content: `
config {
    type: "table",
    schema: "electric_cars",
    dependencies: 'ALL_EV_CARS_DATA',
    bigquery: {
        partitionBy: "MODEL",
        requirePartitionFilter : true,
        clusterBy: ["CITY", "STATE"]
    },
    tags: ["TAG_1"]
}

SELECT * FROM electric_cars WHERE model = $1;

pre_operations {
  ${when(incremental(), ` + "`" + `DELETE
  FROM
    ${self()}
  WHERE
    DATE(PIPELINE_RUN_DATETIME) = CURRENT_DATE()` + "`" + `)}
}

`,
			expected: sqlxFileMetaData{
				numLines:        23,
				configStartLine: 2,
				configEndLine:   12,
				configString: `config {
    type: "table",
    schema: "electric_cars",
    dependencies: 'ALL_EV_CARS_DATA',
    bigquery: {
        partitionBy: "MODEL",
        requirePartitionFilter : true,
        clusterBy: ["CITY", "STATE"]
    },
    tags: ["TAG_1"]
}
`,
				preOperationsStartLine: 16,
				preOperationsEndLine:   22,
				preOperationsString: `pre_operations {
  ${when(incremental(), ` + "`" + `DELETE
  FROM
    ${self()}
  WHERE
    DATE(PIPELINE_RUN_DATETIME) = CURRENT_DATE()` + "`" + `)}
}
`,
				queryString: `SELECT * FROM electric_cars WHERE model = $1;`,
			},
			wantErr: false,
		},

		{
			name: "Minimal config and longer query and comment before config",
			content: `-- some comment
config {
    type: "table",
    schema: "electric_cars"
}


WITH CTE1 AS (
  SELECT
    MAKE
    , COUNTY
    , CITY
    , STATE
    , POSTAL_CODE
    , MODEL
    , MODEL_YEAR
    , COUNT(VIN) AS CNT_VIN
  FROM ${ref("ALL_EV_CARS_DATA")}
  GROUP BY MAKE, COUNTY, CITY, STATE, POSTAL_CODE, MODEL, MODEL_YEAR
  HAVING MAKE = ${constants.make}
)
SELECT * FROM CTE1
            `,
			expected: sqlxFileMetaData{
				numLines:        22,
				configStartLine: 2,
				configEndLine:   5,
				configString: `config {
    type: "table",
    schema: "electric_cars"
}
`,
				queryString: `
WITH CTE1 AS (
  SELECT
    MAKE
    , COUNTY
    , CITY
    , STATE
    , POSTAL_CODE
    , MODEL
    , MODEL_YEAR
    , COUNT(VIN) AS CNT_VIN
  FROM ${ref("ALL_EV_CARS_DATA")}
  GROUP BY MAKE, COUNTY, CITY, STATE, POSTAL_CODE, MODEL, MODEL_YEAR
  HAVING MAKE = ${constants.make}
)
SELECT * FROM CTE1
                `,
			},
			wantErr: false,
		},
		//TODO: Need to handle case where file does not have a config block
		// {
		// 	name: "File without config",
		// 	content: `-- name: GetElectricCars :many
		// SELECT * FROM electric_cars WHERE model = $1;`,
		// 	expected: sqlxFileMetaData{
		// 		numLines:        1,
		// 		configStartLine: 0,
		// 		configEndLine:   0,
		// 		configString:    "",
		// 		queryString: `SELECT * FROM electric_cars WHERE model = $1;`,
		// 	},
		// 	wantErr: false,
		// },
		{
			name:     "Empty file",
			content:  "",
			expected: sqlxFileMetaData{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tmpfile, err := os.CreateTemp("", "test*.sqlx")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			// Write content to the file
			if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}

			// Call the function
			got, err := getSqlxFileMetaData(tmpfile.Name())

			// Check for errors
			if (err != nil) != tt.wantErr {
				t.Errorf("getSqlxFileMetaData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Set the filepath in the expected result
			tt.expected.filepath = tmpfile.Name()

			// Compare each field separately
			if got.filepath != tt.expected.filepath {
				t.Errorf("[got]:  filepath = %v, [want]:  %v", got.filepath, tt.expected.filepath)
			}
			if got.numLines != tt.expected.numLines {
				t.Errorf("[got]:  numLines = %v, [want]:  %v", got.numLines, tt.expected.numLines)
			}
			if got.configStartLine != tt.expected.configStartLine {
				t.Errorf("[got]:  configStartLine = %v, [want]:  %v", got.configStartLine, tt.expected.configStartLine)
			}
			if got.configEndLine != tt.expected.configEndLine {
				t.Errorf("[got]:  configEndLine = %v, [want]:  %v", got.configEndLine, tt.expected.configEndLine)
			}
			if strings.TrimSpace(got.configString) != strings.TrimSpace(tt.expected.configString) {
				t.Errorf("[got]:  configString = %v, [want]:  %v", got.configString, tt.expected.configString)
			}
			if got.preOperationsStartLine != tt.expected.preOperationsStartLine {
				t.Errorf("[got]:  preOperationsStartLine = %v, [want]:  %v", got.preOperationsStartLine, tt.expected.preOperationsStartLine)
			}
			if got.preOperationsEndLine != tt.expected.preOperationsEndLine {
				t.Errorf("[got]:  preOperationsEndLine = %v, [want]:  %v", got.preOperationsEndLine, tt.expected.preOperationsEndLine)
			}
			if strings.TrimSpace(got.preOperationsString) != strings.TrimSpace(tt.expected.preOperationsString) {
				t.Errorf("[got]:  preOperationsString = %v, [want]:  %v", got.preOperationsString, tt.expected.preOperationsString)
			}
			if strings.TrimSpace(got.queryString) != strings.TrimSpace(tt.expected.queryString) {
				t.Errorf("[got]:  queryString = %v, [want]:  %v", got.queryString, tt.expected.queryString)
			}
		})
	}
}

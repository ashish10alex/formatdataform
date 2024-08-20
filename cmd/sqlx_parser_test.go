package cmd

import (
	"os"
	"strings"
	"testing"
)

var simplePostOpsBlock = `
post_operations {
    select 1
    union all
    select 2
}
`

var complexPreOpsBlock = `
pre_operations {
  ${when(incremental(),` +
  "`" + `
  DELETE
  FROM
    ${self()}
  WHERE
    DATE(SNAPSHOT_DATE) = CURRENT_DATE()` +
  "`" + `
    )
  }
}
`

func TestSqlxParser(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected sqlxParserMeta
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
			expected: sqlxParserMeta{
				numLines: 13,
				configBlockMeta: ConfigBlockMeta{
					exsists:            true,
					startOfConfigBlock: 2,
					endOfConfigBlock:   12,
					configBlockContent: `config {
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
`},
				sqlBlocksMeta: SqlBlockMeta{
					exsists:                  true,
					startOfSqlBlock:          13,
					endOfSqlBlock:            13,
					sqlBlockContent:          `SELECT * FROM electric_cars WHERE model = $1;`,
					formattedSqlBlockContent: "",
				},
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
}` + complexPreOpsBlock + simplePostOpsBlock + `
SELECT * FROM electric_cars WHERE model = $1
limit 100
            `,
			expected: sqlxParserMeta{
				numLines: 32,
				configBlockMeta: ConfigBlockMeta{
					exsists:            true,
					startOfConfigBlock: 2,
					endOfConfigBlock:   12,
					configBlockContent: `
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
}`,
				},
				sqlBlocksMeta: SqlBlockMeta{
					exsists:         true,
					startOfSqlBlock: 30,
					endOfSqlBlock:   32,
					sqlBlockContent: `SELECT * FROM electric_cars WHERE model = $1
limit 100
                    `,
					formattedSqlBlockContent: "",
				},
				preOpsBlocksMeta: []PreOpsBlockMeta{
					{
						exsists:                   true,
						startOfPreOperationsBlock: 13,
						endOfPreOperationsBlock:   22,
						preOpsBlockContent:        strings.TrimPrefix(complexPreOpsBlock, "\n"),
					},
				},
				postOpsBlocksMeta: []PostOpsBlockMeta{
					{
						exsists:                    true,
						startOfpostOperationsBlock: 24,
						endOfpostOperationsBlock:   28,
						postOpsBlockContent:        strings.TrimPrefix(simplePostOpsBlock, "\n"),
					},
				},
			},
			wantErr: false,
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
			got, err := sqlxParser(tmpfile.Name())

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
			if got.configBlockMeta.startOfConfigBlock != tt.expected.configBlockMeta.startOfConfigBlock {
				t.Errorf("[got]:  configStartLine = %v, [want]:  %v", got.configBlockMeta.startOfConfigBlock, tt.expected.configBlockMeta.startOfConfigBlock)
			}
			if got.configBlockMeta.endOfConfigBlock != tt.expected.configBlockMeta.endOfConfigBlock {
				t.Errorf("[got]:  configEndLine = %v, [want]:  %v", got.configBlockMeta.endOfConfigBlock, tt.expected.configBlockMeta.endOfConfigBlock)
			}
			if strings.TrimSpace(got.configBlockMeta.configBlockContent) != strings.TrimSpace(tt.expected.configBlockMeta.configBlockContent) {
				t.Errorf("[got]:  configString = %v, [want]:  %v", got.configBlockMeta.configBlockContent, tt.expected.configBlockMeta.configBlockContent)
			}
			if strings.TrimSpace(got.sqlBlocksMeta.sqlBlockContent) != strings.TrimSpace(tt.expected.sqlBlocksMeta.sqlBlockContent) {
				t.Errorf("[got]:  sqlBlockContent = %v, [want]:  %v", got.sqlBlocksMeta.sqlBlockContent, tt.expected.sqlBlocksMeta.sqlBlockContent)
			}

			if got.sqlBlocksMeta.startOfSqlBlock != tt.expected.sqlBlocksMeta.startOfSqlBlock {
				t.Errorf("[got]:  startOfSqlBlock = %v, [want]:  %v", got.sqlBlocksMeta.startOfSqlBlock, tt.expected.sqlBlocksMeta.startOfSqlBlock)
			}
			if got.sqlBlocksMeta.endOfSqlBlock != tt.expected.sqlBlocksMeta.endOfSqlBlock {
				t.Errorf("[got]:  endOfSqlBlock = %v, [want]:  %v", got.sqlBlocksMeta.endOfSqlBlock, tt.expected.sqlBlocksMeta.endOfSqlBlock)
			}

			if len(tt.expected.preOpsBlocksMeta) > 0 || len(tt.expected.preOpsBlocksMeta) > 0 {

				if got.preOpsBlocksMeta[0].startOfPreOperationsBlock != tt.expected.preOpsBlocksMeta[0].startOfPreOperationsBlock {
					t.Errorf("[got]:  startOfPreOperationsBlock = %v, [want]:  %v", got.preOpsBlocksMeta[0].startOfPreOperationsBlock, tt.expected.preOpsBlocksMeta[0].startOfPreOperationsBlock)
				}

				if got.preOpsBlocksMeta[0].endOfPreOperationsBlock != tt.expected.preOpsBlocksMeta[0].endOfPreOperationsBlock {
					t.Errorf("[got]:  endOfPreOperationsBlock = %v, [want]:  %v", got.preOpsBlocksMeta[0].endOfPreOperationsBlock, tt.expected.preOpsBlocksMeta[0].endOfPreOperationsBlock)
				}

				if got.preOpsBlocksMeta[0].preOpsBlockContent != tt.expected.preOpsBlocksMeta[0].preOpsBlockContent {
					t.Errorf("[got]:  preOpsBlockContent = %q, [want]:  %q", got.preOpsBlocksMeta[0].preOpsBlockContent, tt.expected.preOpsBlocksMeta[0].preOpsBlockContent)
				}

			}

			if len(tt.expected.postOpsBlocksMeta) > 0 || len(tt.expected.postOpsBlocksMeta) > 0 {
				if got.postOpsBlocksMeta[0].startOfpostOperationsBlock != tt.expected.postOpsBlocksMeta[0].startOfpostOperationsBlock {
					t.Errorf("[got]:  startOfpostOperationsBlock = %v, [want]:  %v", got.postOpsBlocksMeta[0].startOfpostOperationsBlock, tt.expected.postOpsBlocksMeta[0].startOfpostOperationsBlock)
				}

				if got.postOpsBlocksMeta[0].endOfpostOperationsBlock != tt.expected.postOpsBlocksMeta[0].endOfpostOperationsBlock {
					t.Errorf("[got]:  endOfpostOperationsBlock = %v, [want]:  %v", got.postOpsBlocksMeta[0].endOfpostOperationsBlock, tt.expected.postOpsBlocksMeta[0].endOfpostOperationsBlock)
				}
				if got.postOpsBlocksMeta[0].postOpsBlockContent != tt.expected.postOpsBlocksMeta[0].postOpsBlockContent {
					t.Errorf("[got]:  postOpsBlockContent = %q, [want]:  %q", got.postOpsBlocksMeta[0].postOpsBlockContent, tt.expected.postOpsBlocksMeta[0].postOpsBlockContent)
				}
			}

		})
	}
}

package cmd


var pythonCode = `
import sys
sqlfluff_config_path = sys.argv[1]
my_bad_query = sys.argv[2]

def fix_query(sqlfluff_config_path, my_bad_query):
    try:
        import sqlfluff
    except ImportError:
        return "sqlfluff is not installed"
    my_good_query = sqlfluff.fix(
                        my_bad_query,
                        dialect="bigquery",
                        config_path=str(sqlfluff_config_path),
                        )
    return my_good_query

print(fix_query(sqlfluff_config_path, my_bad_query)) # so that GoLang can read stdout
        `

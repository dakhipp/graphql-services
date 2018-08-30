/*====
Terraform config
======*/

bastion_public_key = "### Too long no example ###"

github_oauth = "### oatuh token ###"

github_user = "example-user"

github_repo = "example-repo"

ssl_identifier = "example.com"

route53_zone_id = "ABC1234AB12A1"

/*====
Terraform and migration container config
======*/

psql_root_db = "psql"

psql_root_user = "root"

psql_root_pass = "toor1234"

/*====
Terraform config, must be unique across stages
======*/

environment = "dev"

github_branch = "master"

domain = "example.example.com"

artifact_bucket_name = "artifcate-bucket-s3"

/*====
Container Env config
======*/

psql_db = "example_db"

psql_user = "example_user"

psql_pass = "example_pass"

psql_port = 5432

graphql_port = 8000

auth_port = 8001

playground_enabled = true

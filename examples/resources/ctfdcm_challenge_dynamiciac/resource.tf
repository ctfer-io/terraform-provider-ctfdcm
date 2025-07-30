resource "ctfdcm_challenge_dynamiciac" "http" {
  name        = "My Challenge"
  category    = "misc"
  description = "..."
  value       = 500
  decay       = 100
  minimum     = 50
  state       = "visible"
  function    = "logarithmic"

  shared          = true
  destroy_on_flag = true
  mana_cost       = 1
  scenario        = "localhost:5000/some/scenario:v0.1.0"
  timeout         = 600

  topics = [
    "Misc"
  ]
  tags = [
    "misc",
    "basic"
  ]
}

resource "ctfd_flag" "http_flag" {
  challenge_id = ctfdcm_challenge_dynamiciac.http.id
  content      = "CTF{some_flag}"
}

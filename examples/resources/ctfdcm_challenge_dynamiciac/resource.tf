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
  scenario_id     = ctfd_file.scenario.id
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

resource "ctfd_file" "scenario" {
  name       = "scenario.zip"
  contentb64 = filebase64(".../scenario.zip")
}

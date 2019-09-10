variable "region" {
  default = "ap-southeast-2"
}

variable "function-local-directory" {
  default = ""
  description = "Work around for calculating hash of file as can't use version numbers"
}

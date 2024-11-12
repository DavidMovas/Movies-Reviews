variable "instance_type" {
  type = string
  default = "t3.micro"
  description = "The type of EC2 instance"
}

variable "ami_id" {
  type = string
  default = "ami-00ac244ee0ad9050d"
}

variable "AppName" {
  type = string
  default = "movie-reviews"
  description = "The name of the application"
}

variable "secrets" {
  type = set(string)
  default = [
    "/movie-reviews/jwt-secret",
    "/movie-reviews/admin/name",
    "/movie-reviews/admin/email",
    "/movie-reviews/admin/password"
  ]
}
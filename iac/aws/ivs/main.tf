# main.tf
provider "aws" {
  region = "us-east-1"
}

resource "aws_ivs_channel" "this" {
  name         = "majesticcoding"
  latency_mode = "LOW"
  type         = "STANDARD"

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_ivs_stream_key" "this" {
  channel_arn = aws_ivs_channel.this.arn
}

output "rtmp_server_url" {
  value = "rtmps://${aws_ivs_channel.this.ingest_endpoint}:443/app/"
}

output "playback_url" {
  value = aws_ivs_channel.this.playback_url
}

output "stream_key" {
  value     = aws_ivs_stream_key.this.value
  sensitive = true
}

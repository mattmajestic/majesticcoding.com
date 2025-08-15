package models

type CloudBrand struct {
	Name string
	Logo string // URL or local path to image
	Alt  string
	Link string // optional, if you want the brand clickable
}

var CloudBrands = []CloudBrand{
	{
		Name: "Microsoft Azure",
		Logo: "https://upload.wikimedia.org/wikipedia/commons/f/fa/Microsoft_Azure.svg",
		Alt:  "Azure Logo",
		Link: "https://azure.microsoft.com/",
	},
	{
		Name: "Amazon Web Services",
		Logo: "https://upload.wikimedia.org/wikipedia/commons/9/93/Amazon_Web_Services_Logo.svg",
		Alt:  "AWS Logo",
		Link: "https://aws.amazon.com/",
	},
	{
		Name: "Google Cloud Platform",
		Logo: "https://upload.wikimedia.org/wikipedia/commons/5/51/Google_Cloud_logo.svg",
		Alt:  "GCP Logo",
		Link: "https://cloud.google.com/",
	},
}

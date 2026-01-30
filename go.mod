module smtp-sender

go 1.25.4

require (
	github.com/joho/godotenv v1.5.1
	github.com/xuri/excelize/v2 v2.10.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	pay_slip_generator v0.0.0-00010101000000-000000000000
)

require (
	github.com/jung-kurt/gofpdf v1.16.2 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/tiendc/go-deepcopy v1.7.1 // indirect
	github.com/xuri/efp v0.0.1 // indirect
	github.com/xuri/nfp v0.0.2-0.20250530014748-2ddeb826f9a9 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)

replace pay_slip_generator => ./pay_slip_generator/pay_slip_generator

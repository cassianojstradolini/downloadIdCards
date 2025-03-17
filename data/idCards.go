package data

const (
	mockMedicalIdCardFrontURL = "https://images.ctfassets.net/ii2mrg7wcwxq/5ulus6CP80PUtJTBkrsPmH/67ca543b65c527fd2a2530bb05ed16f2/med_card_front.png"
	mockMedicalIdCardBackURL  = "https://images.ctfassets.net/ii2mrg7wcwxq/7p2j22SpflOacalSecIrCV/1a14e8f2cfd7fa5bb5f253d61e97f5db/med_card_back.png"
)

const (
	mockIdCardFrontFaceAltText = `
		HEALTHCO
		Health Plan (80840) 999-12345-99
		Member ID: 123456789 Group Number: 123456
		Member: SAMPLE A SAMPLE Payer ID: 12345
		PCP Name: DR SAMPLE M SAMPLE MD Rx Bin: 999999
		PCP Phone: 999-999-9999 Rx Grp: 999999 Rx PCN: 9999
		0501 Healthcare Community Plan
		Administered by [Appropriate legal entity]
	`
	mockIdCardBackFaceAltText = `
		In an emergency go to nearest emergency room or call 911. Printed: 01/30/23
		By using this card for services, you agree to the release of medical information, as stated in your member handbook. To verify benefits or to find a provider, visit the website
		For Members: 999-999-9999 TTY 711
		NurseLine: 999-999-9999 TTY 711
		For Providers: HCprovider.com 999-999-9999
		Medical Claims: PO Box 9999, Kingston, NY, 99999-8207
		Pharmacy Claims: RX, PO Box 999999, Dallas, TX 99999-0334
		For Pharmacists: 999-999-9999
	`
)

var MockIdCardFront = IdCard{
	Id:   "mock-id-card-front",
	Type: IdCardTypeIdCard,
	Attributes: IdCardAttributes{
		BenefitId:   nil,
		BenefitType: nil,
		Type:        IdCardAttributesTypeUrl,
		Face:        IdCardAttributesFaceFront,
		Source:      mockMedicalIdCardFrontURL,
		AltText:     mockIdCardFrontFaceAltText,
	},
}

var MockIdCardBack = IdCard{
	Id:   "mock-id-card-back",
	Type: IdCardTypeIdCard,
	Attributes: IdCardAttributes{
		BenefitId:   nil,
		BenefitType: nil,
		Type:        IdCardAttributesTypeUrl,
		Face:        IdCardAttributesFaceBack,
		Source:      mockMedicalIdCardBackURL,
		AltText:     mockIdCardBackFaceAltText,
	},
}

var MockImageIdCardFront = IdCard{
	Id:   "mock-image-id-card-front",
	Type: IdCardTypeIdCard,
	Attributes: IdCardAttributes{
		BenefitId:   nil,
		BenefitType: nil,
		Type:        IdCardAttributesTypeBase64,
		Face:        IdCardAttributesFaceFront,
		Source:      frontBase64Src,
		AltText:     mockIdCardFrontFaceAltText,
	},
}

var MockImageIdCardBack = IdCard{
	Id:   "mock-image-id-card-back",
	Type: IdCardTypeIdCard,
	Attributes: IdCardAttributes{
		BenefitId:   nil,
		BenefitType: nil,
		Type:        IdCardAttributesTypeBase64,
		Face:        IdCardAttributesFaceBack,
		Source:      backBase64Src,
		AltText:     mockIdCardFrontFaceAltText,
	},
}

var MockHTMLIdCardFront = IdCard{
	Id:   "mock-html-id-card-front",
	Type: IdCardTypeIdCard,
	Attributes: IdCardAttributes{
		BenefitId:   nil,
		BenefitType: nil,
		Type:        IdCardAttributesTypeHTML,
		Face:        IdCardAttributesFaceFront,
		Source:      htmlIdCardFront,
		AltText:     mockIdCardFrontFaceAltText,
	},
}

var MockHTMLIdCardBack = IdCard{
	Id:   "mock-html-id-card-back",
	Type: IdCardTypeIdCard,
	Attributes: IdCardAttributes{
		BenefitId:   nil,
		BenefitType: nil,
		Type:        IdCardAttributesTypeHTML,
		Face:        IdCardAttributesFaceBack,
		Source:      htmlIdCardBack,
	},
}

var MockHTMLIdCardBoth = IdCard{
	Id:   "mock-html-id-card-both",
	Type: IdCardTypeIdCard,
	Attributes: IdCardAttributes{
		BenefitId:   nil,
		BenefitType: nil,
		Type:        IdCardAttributesTypeHTML,
		Face:        "combined",
		Source:      htmlIdCardBoth,
	},
}

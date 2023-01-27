package seeder

import (
	"dbseeder/internal/schema"
	"errors"
	"strconv"
	"strings"
	"syreclabs.com/go/faker"
	"time"
)

const (
	fieldString         schema.FieldType = "string"
	fieldInt            schema.FieldType = "int"
	fieldEmail          schema.FieldType = "email"
	fieldMac            schema.FieldType = "mac"
	fieldIpv4           schema.FieldType = "ipv4"
	fieldIpv6           schema.FieldType = "ipv6"
	fieldDate           schema.FieldType = "date"
	fieldText           schema.FieldType = "text"
	fieldAddress        schema.FieldType = "address"
	fieldAddressZipCode schema.FieldType = "address->zip"
	fieldAddressCity    schema.FieldType = "address->city"
	fieldAddressStreet  schema.FieldType = "address->street"
	fieldCountry        schema.FieldType = "country"
	fieldCountryCode    schema.FieldType = "country->code"
	fieldMoney          schema.FieldType = "money"
	fieldUrl            schema.FieldType = "url"
	fieldDomainName     schema.FieldType = "domain"
	fieldFullName       schema.FieldType = "name"
	fieldFirstName      schema.FieldType = "firstName"
	fieldLastName       schema.FieldType = "fieldLastName"
	fieldHex            schema.FieldType = "hex"
	fieldDecimal        schema.FieldType = "decimal"
	fieldPhoneNumber    schema.FieldType = "phone"
	fieldPhoneCode      schema.FieldType = "phone->code"
)

var FieldTypesMap = map[schema.FieldType]string{
	fieldString:         "String field can add length like - 'string 15'",
	fieldInt:            "Int32 field. Min Max - 'int 0 10'",
	fieldEmail:          "Email field",
	fieldMac:            "Mac Address",
	fieldIpv4:           "IPv4 address",
	fieldIpv6:           "IPv6 address",
	fieldDate:           "Date field. Y-m-d format",
	fieldText:           "Random text field like 'Lorem Ipsam'",
	fieldAddress:        "Full address field",
	fieldAddressZipCode: "Address zipcode",
	fieldAddressCity:    "Address city",
	fieldAddressStreet:  "Address street",
	fieldCountry:        "Country",
	fieldCountryCode:    "Country code",
	fieldMoney:          "Money",
	fieldUrl:            "Url",
	fieldDomainName:     "Domain name",
	fieldFullName:       "Field full name",
	fieldFirstName:      "Field first name",
	fieldLastName:       "Field last name",
	fieldHex:            "Field HEX",
	fieldDecimal:        "Decimal",
	fieldPhoneNumber:    "Phone number",
	fieldPhoneCode:      "Phone country code",
}

type Modifier interface {
	Apply(modifierName string, val any) (any, error)
}

func Generate(t string) (any, error) {
	parts := strings.Split(t, " ")
	if len(parts) <= 0 {
		return nil, errors.New("wrong type format")
	}

	fieldType := schema.FieldType(parts[0])
	switch fieldType {
	case fieldString:
		size := 10
		if len(parts) <= 1 {
			size, _ = strconv.Atoi(parts[1])
		}
		return faker.Lorem().Characters(size), nil
	case fieldInt:
		begin, end := 0, 1000
		if len(parts) == 3 {
			begin, _ = strconv.Atoi(parts[1])
			end, _ = strconv.Atoi(parts[2])
		}
		if len(parts) == 2 {
			end, _ = strconv.Atoi(parts[1])
		}
		return faker.RandomInt(begin, end), nil
	case fieldMac:
		return faker.Internet().MacAddress(), nil
	case fieldIpv4:
		return faker.Internet().IpV4Address(), nil
	case fieldIpv6:
		return faker.Internet().IpV6Address(), nil
	case fieldDecimal:
		begin, end := 0, 1000
		if len(parts) == 3 {
			begin, _ = strconv.Atoi(parts[1])
			end, _ = strconv.Atoi(parts[2])
		}
		if len(parts) == 2 {
			end, _ = strconv.Atoi(parts[1])
		}
		return faker.Number().Decimal(begin, end), nil
	case fieldHex:
		size := 10
		if len(parts) <= 1 {
			size, _ = strconv.Atoi(parts[1])
		}
		return faker.Number().Hexadecimal(size), nil
	case fieldEmail:
		return faker.Internet().SafeEmail(), nil
	case fieldPhoneCode:
		return faker.PhoneNumber().AreaCode(), nil
	case fieldPhoneNumber:
		return faker.PhoneNumber().PhoneNumber(), nil
	case fieldFullName:
		return faker.Name().Name(), nil
	case fieldFirstName:
		return faker.Name().FirstName(), nil
	case fieldLastName:
		return faker.Name().LastName(), nil
	case fieldUrl:
		return faker.Internet().Url(), nil
	case fieldDomainName:
		return faker.Internet().DomainName(), nil
	case fieldText:
		return faker.Lorem().Paragraph(4), nil
	case fieldDate:
		begin, _ := time.Parse("2006-01-02", "1990-01-01")
		end := time.Now()
		if len(parts) == 3 {
			begin, _ = time.Parse("2006-01-02", parts[1])
			end, _ = time.Parse("2006-01-02", parts[2])
		}
		if len(parts) == 2 {
			begin, _ = time.Parse("2006-01-02", parts[1])
		}
		return faker.Date().Between(begin, end), nil
	case fieldMoney:
		return faker.Commerce().Price(), nil
	case fieldAddress:
		return faker.Address().StreetAddress(), nil
	case fieldCountry:
		return faker.Address().Country(), nil
	case fieldCountryCode:
		return faker.Address().CountryCode(), nil
	case fieldAddressCity:
		return faker.Address().City(), nil
	case fieldAddressStreet:
		return faker.Address().StreetName(), nil
	case fieldAddressZipCode:
		return faker.Address().ZipCode(), nil
	default:
		return nil, errors.New("unknown field type")
	}
}

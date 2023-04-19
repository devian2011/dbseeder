package fake

import (
	"dbseeder/internal/schema"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"syreclabs.com/go/faker"
	"time"
)

type FieldType string

const (
	fieldNull           FieldType = "null"
	fieldEmpty          FieldType = "empty"
	fieldBool           FieldType = "bool"
	fieldString         FieldType = "string"
	fieldInt            FieldType = "int"
	fieldEmail          FieldType = "email"
	fieldMac            FieldType = "mac"
	fieldIpv4           FieldType = "ipv4"
	fieldIpv6           FieldType = "ipv6"
	fieldDate           FieldType = "date"
	fieldText           FieldType = "text"
	fieldFullAddress    FieldType = "address->full"
	fieldAddress        FieldType = "address"
	fieldAddressZipCode FieldType = "address->zip"
	fieldAddressCity    FieldType = "address->city"
	fieldAddressStreet  FieldType = "address->street"
	fieldCountry        FieldType = "country"
	fieldCountryCode    FieldType = "country->code"
	fieldMoney          FieldType = "money"
	fieldURL            FieldType = "url"
	fieldDomainName     FieldType = "domain"
	fieldFullName       FieldType = "name"
	fieldFirstName      FieldType = "firstName"
	fieldLastName       FieldType = "lastName"
	fieldHex            FieldType = "hex"
	fieldDecimal        FieldType = "decimal"
	fieldPhoneNumber    FieldType = "phone"
	fieldPhoneCode      FieldType = "phone->code"
)

var FieldTypesMap = map[FieldType]string{
	fieldNull:           "Null field",
	fieldEmpty:          "Empty string field",
	fieldBool:           "Boolean field",
	fieldString:         "String field can add length like - 'string 15'",
	fieldInt:            "Int32 field. Min Max - 'int 0 10'",
	fieldEmail:          "Email field",
	fieldMac:            "Mac Address",
	fieldIpv4:           "IPv4 address",
	fieldIpv6:           "IPv6 address",
	fieldDate:           "Date field. Y-m-d format",
	fieldText:           "Random text field like 'Lorem Ipsam'",
	fieldFullAddress:    "Full address field",
	fieldAddress:        "Full street address field",
	fieldAddressZipCode: "Address zipcode",
	fieldAddressCity:    "Address city",
	fieldAddressStreet:  "Address street",
	fieldCountry:        "Country",
	fieldCountryCode:    "Country code",
	fieldMoney:          "Money",
	fieldURL:            "Url",
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

func Generate(fieldName string, fieldVal schema.Field) (any, error) {
	parts := strings.Split(fieldVal.Type, " ")
	if len(parts) <= 0 {
		return nil, fmt.Errorf("wrong type format for %s", fieldName)
	}

	fieldType := FieldType(parts[0])
	switch fieldType {
	case fieldNull:
		return nil, nil
	case fieldEmpty:
		return "", nil
	case fieldBool:
		rand.Seed(time.Now().UnixNano())
		return rand.Intn(2) == 1, nil
	case fieldString:
		size := 10
		if len(parts) >= 2 {
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
		if len(parts) >= 2 {
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
	case fieldURL:
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
	case fieldFullAddress:
		return faker.Address().String(), nil
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
		return nil, fmt.Errorf("unknown field type for %s", fieldName)
	}
}

package mailparser

import (
	"io/ioutil"
	"testing"
)

func TestEmailWithQuotingParsing(t *testing.T) {
	data, err := ioutil.ReadFile("email_test_quoted.eml")
	if err != nil {
		t.Fatal(err)
	}

	email, err := ParseEmail(data)
	if err != nil {
		t.Fatal(err)
	}

	if email.Header("Subject") != "[USN-1108-2] DHCP vulnerability" {
		t.Fatal("Cannot read the Subject properly")
	}
	//	fmt.Println(email.GetHeader("Content-Transfer-Encoding"))

	if len(email.Body()) == 0 {
		t.Fatal("Cannot parse the quoted email body correctly")
	}

	if len(email.parts) != 2 {
		t.Fatal("Failed to parse the email message parts")
	}
}

func TestEmailParsing(t *testing.T) {
	data, err := ioutil.ReadFile("email_test_7bit.eml")
	if err != nil {
		t.Fatal(err)
	}

	email, err := ParseEmail(data)
	if err != nil {
		t.Fatal(err)
	}

	if email.Header("Subject") != "[USN-2706-1] OpenJDK 6 vulnerabilities" {
		t.Fatal("Cannot read the Subject properly")
	}
	//	fmt.Println(email.GetHeader("Content-Transfer-Encoding"))

	if len(email.Body()) == 0 {
		t.Fatal("Cannot parse the email body correctly")
	}

	if len(email.parts) != 2 {
		t.Fatal("Failed to parse the email message parts")
	}

}

func TestEmailParsingMultiPartSigned(t *testing.T) {
	data, err := ioutil.ReadFile("email_multi_part_signed.eml")
	if err != nil {
		t.Fatal(err)
	}

	email, err := ParseEmail(data)
	if err != nil {
		t.Fatal(err)
	}

	if email.Header("Subject") != "[USN-1436-1] Libtasn1 vulnerability" {
		t.Fatal("Cannot read the Subject properly")
	}
	//	fmt.Println(email.GetHeader("Content-Transfer-Encoding"))

	if len(email.Body()) == 0 {
		t.Fatal("Cannot parse the email body correctly")
	}

	if len(email.parts) != 2 {
		t.Fatal("Failed to parse the email message parts")
	}

}

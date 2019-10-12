package udwDnsPacket

func IsInMatchOut(in DnsPacket, out DnsPacket) bool {
	return in.Id == out.Id &&
		in.IsRecursionDesired() == out.IsRecursionDesired() &&
		in.QuestionName == out.QuestionName &&
		in.QuestionType == out.QuestionType &&
		in.IsResponse() == false && out.IsResponse() == true
}
func IsLookLikeFakePacket(out DnsPacket) bool {

	return (len(out.AnswerList) == 1) &&
		out.Bits == 0x8180 &&
		out.IsRecursionAvailable() &&
		out.IsRecursionDesired() &&
		out.IsAuthoritative() == false &&
		out.IsResponse() == true &&
		out.HasEdns4096 == false &&
		out.HasCompressDomain == true &&
		out.Ancount == 1 &&
		out.Nscount == 0 &&
		out.Arcount == 0 &&
		out.GetRcode() == RcodeNOERROR &&
		out.GetOpcode() == 0 &&
		out.AnswerList[0].Type == DnsTypeA &&
		out.QuestionName == out.AnswerList[0].Name
}

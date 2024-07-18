package parser

type Instruction struct {
	Result    *Write    `(@@ "=")?`
	Operation Operation `@@ ";"`
}

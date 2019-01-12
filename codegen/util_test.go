package codegen_test

//go:generate mockgen -package=codegen -destination=./mock_type_info.go github.com/lonegunmanb/varys/ast TypeInfo
//go:generate mockgen -package=codegen -destination=./mock_field_info.go github.com/lonegunmanb/varys/ast FieldInfo
//go:generate mockgen -package=codegen -destination=./mock_embedded_type.go github.com/lonegunmanb/varys/ast EmbeddedType
//go:generate mockgen -source=./assembler.go -package=codegen -destination=./mock_assembler.go
//go:generate mockgen -source=./product_type_info_wrap.go -package=codegen -destination=./mock_type_codegen.go
//go:generate mockgen -package=codegen -destination=./mock_gopathenv.go github.com/lonegunmanb/varys/ast GoPathEnv
//go:generate mockgen -source=./register_code_writer.go -package=codegen -destination=./mock_register.go
//go:generate mockgen -source=./pkg_name_arbitrator.go -package=codegen -destination=./mock_pkg_name_arbitrator.go

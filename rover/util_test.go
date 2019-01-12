package rover_test

//go:generate mockgen -package=rover -destination=./mock_gopathenv.go github.com/lonegunmanb/varys/ast GoPathEnv
//go:generate mockgen -package=rover -destination=./mock_file_retriever.go github.com/lonegunmanb/varys/ast FileRetriever
//go:generate mockgen -package=rover -destination=./mock_file_info.go github.com/lonegunmanb/varys/ast FileInfo
//go:generate mockgen -package=rover -destination=./mock_type_walker.go github.com/lonegunmanb/varys/ast TypeWalker

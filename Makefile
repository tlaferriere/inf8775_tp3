########################################
# VARIABLES
########################################
EXAMPLE_SIZES := 100 500 1000 5000 10000 50000 100000
NUM_EXAMPlES=10
EXAMPLE_NUMS := 0 1 2 3 4 5 6 7 8 9
ALGOS := glouton progdyn tabou

# directories
DATA_DIR=data

# files
# Go files
GO_BIN=tp3
GO_SRCS := $(wildcard *.go)

# python
PYTHON=python3

# Example files
EXAMPLE_PREFIX := $(if $(DATA_DIR),$(DATA_DIR)/)N
EXAMPLE_OUT_PREFIX := $(if $(DATA_DIR),$(DATA_DIR)/)O

########################################
# BUILD BIN
########################################
.PHONY: all
all: tp2

########################################
# CREATE REMISE ARCHIVE
########################################
.PHONY: remise
remise: 1905759.zip

%.zip: $(GO_BIN) TP3_H22_Rapport.docx $(GO_SRCS) tp.sh go.mod go.sum Makefile rapport.pdf
	zip -u $@ $^


%.pdf: %.ipynb
	pandoc -s -o $@ $^

########################################
# BUILD GO BINARY
########################################
$(GO_BIN): $(GO_SRCS) generate-pb-go
	go build -o $(GO_BIN)

########################################
# BUILD GO BINARY
########################################
ALL_EX=$(foreach x,$(EXAMPLE_NUMS),verify-100_$(x)_$(a))
VERIFY_ALL_100 := $(foreach a,$(ALGOS),$(ALL_EX))
test: $(VERIFY_ALL_100)

verify-%: $(EXAMPLE_OUT_PREFIX)%.txt
	$(PYTHON) verify_TP2.py -s $<

$(EXAMPLE_OUT_PREFIX)%_glouton.txt: $(EXAMPLE_PREFIX)%.txt $(GO_BIN)
	./tp.sh -a glouton -e $< -p > $@

$(EXAMPLE_OUT_PREFIX)%_progdyn.txt: $(EXAMPLE_PREFIX)%.txt $(GO_BIN)
	./tp.sh -a progdyn -e $< -p > $@

$(EXAMPLE_OUT_PREFIX)%_tabou.txt: $(EXAMPLE_PREFIX)%.txt $(GO_BIN)
	./tp.sh -a tabou -e $< -p > $@

$(EXAMPLE_PREFIX)%.txt: gen-examples

.PHONY: gen-examples
gen-examples:
	$PYTHON inst_gen.py && $(if $(DATA_DIR),mv N*_*.txt $(DATA_DIR)) -t 100 -k 3 -n 10

########################################
# CLEAN GENERATED FILES
########################################
.PHONY: clean
clean:
	rm -f $(GO_BIN)
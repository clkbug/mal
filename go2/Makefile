BINS = step0_repl \
  step1_read_print \
  step2_eval \
  step3_env \
  step4_if_fn_do \
  #step5_tco \
  step6_file \
  step7_quote \
  step8_macros \
  step9_try \
  stepA_mal


#####################

all: $(BINS)

dist: mal

mal: $(word $(words $(BINS)),$(BINS))
	cp $< $@

define dep_template
$(1): $(shell ls src/$(1)/*.go)
	cd src/$(1) && go build && mv $(1) ../../
endef

$(foreach b,$(BINS),$(eval $(call dep_template,$(b))))

clean:
	rm -f $(BINS) mal

.PHONY: stats stats-lisp

stats: $(SOURCES)
	@wc $^
	@printf "%5s %5s %5s %s\n" `grep -E "^[[:space:]]*//|^[[:space:]]*$$" $^ | wc` "[comments/blanks]"
stats-lisp: $(SOURCES_LISP)
	@wc $^
	@printf "%5s %5s %5s %s\n" `grep -E "^[[:space:]]*//|^[[:space:]]*$$" $^ | wc` "[comments/blanks]"

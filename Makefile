# Makefile contains targets to assist in building viztransform commands.
#
# Ones that don't build commands log what they are doing followed by a newline.
# Ones that do build commands echo 'making' and what they're making followed by
# a newline.
.PHONY: doc

# all builds the all commands and generates docs.
all: viztransform_apply viztransform_simplify viztransform_viz doc

# viztransform_apply makes the viztransform_apply command.
viztransform_apply:
	@echo "making $@"
	$(call go,$@)
	@echo

# viztransform_simplify makes the viztransform_simplify command.
viztransform_simplify:
	@echo "making $@"
	$(call go,$@)
	@echo

# viztransform_viz makes the viztransform_viz command.
viztransform_viz:
	@echo "making $@"
	$(call go,$@)
	@echo

# doc makes the docs.
doc:
	@echo 'making doc'
	@echo


# pdf is used to make PDFs from Latex files using Pandoc.
#
# The Latex files are expected to be found in the doc directory.
define pdf
	pandoc doc/$(1).md --latex-engine xelatex -o doc/$(1).pdf
endef

# go is used to install Go commands referred to by the name of the command.
#
# The commands are expected to be found in the cmd directory.
define go
	cd cmd/$(1) && go install
endef

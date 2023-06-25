# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

###  Breaking Changes

### Fixed

### Changed

### Added

## [0.2.2]

- Fixed incorrect parsing at `pretrained/createReplaceDecoder()` function

## [0.2.1]

- [#26]: load pretrained Roberta tokenizer failure
- Fixed errors at pretained.FromFile
- upgrade "golang.org/x/text" to fix github warning.

## [0.2.0]

###  Breaking Changes

- `processor.NewRobertaProcessing` added 2 new parameters.
- added `SequenceRanges` field to `Encoding` struct

### Added
- Completed list of pretokenizers, decoders, normalizers, processors

## [0.1.16]
- Fixed data race at `PostProcess` and `EncodeBatch`
- Error handling when `Tokenizer Model` is nil.

## [0.1.12]

### Changed
- Clean-up unwanted console print-out at `processor/bert`

## [0.1.11]

### Added
- Added pretrained model "Roberta" and "GPT2" and "BertLargeCasedWholeWordMaskingSquad"

### Fixed
- [#14]: fixed truncation and padding not working properly

### Changed
- Update "example/truncation", "example/pretrained"

## [0.1.10]

### Fixed
- [#13]: fixed Wordpiece Decoder not join incorrectly tokens and not strip prefix.

## [0.1.9]

## Fixed
- [#13]: fixed Wordpiece Decoder not join incorrectly tokens and not strip prefix.

## [0.1.8]

### Fixed
- [#12]: fixed using pointer to Decoder interface in Tokenizer struct.

## [0.1.7]

### Changed
- Updated `example_test` and `example` in README

### Added
- [#11]: added `addSpecialTokens` param to `EncodeSingle` `EncodePair` and `Tokenizer` APIs.


## [0.1.6]

### Changed
- Update Changelog and README

## [0.1.5]

### Added
- [#10]: setup pretrained subpackage to load pretrained tokenizers. 

## [0.1.4]

### Fixed
- [#8]: Fixed `encoding.MergeWith` merge overflowing incorrectly. 

[#8]: https://github.com/sugarme/tokenizer/pull/8
[#10]: https://github.com/sugarme/tokenizer/pull/10
[#11]: https://github.com/sugarme/tokenizer/issues/11
[#12]: https://github.com/sugarme/tokenizer/issues/12
[#13]: https://github.com/sugarme/tokenizer/issues/13
[#14]: https://github.com/sugarme/tokenizer/issues/14

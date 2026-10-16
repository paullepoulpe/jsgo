const MSG_ICU_EXAMPLE = declareIcuTemplate(
    `{NUM_PEOPLE, plural, offset:1
        =0 {I see no one at all in {INTERPOLATION_2}.}
        =1 {I see {INTERPOLATION_1} in {INTERPOLATION_2}.}
        =2 {I see {INTERPOLATION_1} and one other person in {INTERPOLATION_2}.}
        other {I see {INTERPOLATION_1} and # other people in {INTERPOLATION_2}.}
    }`,
    {
      description: 'ICU example message',
      example: {
        'INTERPOLATION_1': 'Jane Doe',
// This placeholder name doesn't match any placeholders in the message, so it should
// be reported as an error
        'INTERPOLATION_TYPO_2': 'Paris, France',
      }
    });

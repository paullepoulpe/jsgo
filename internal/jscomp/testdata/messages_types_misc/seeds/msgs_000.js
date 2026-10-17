const MSG_ICU_EXAMPLE = declareIcuTemplate(
    `{NUM_PEOPLE, plural, offset:1
        =0 {I see no one at all in {INTERPOLATION_2}.}
        =1 {I see {INTERPOLATION_1} in {INTERPOLATION_2}.}
        =2 {I see {INTERPOLATION_1} and one other person in {INTERPOLATION_2}.}
        other {I see {INTERPOLATION_1} and # other people in {INTERPOLATION_2}.}
    }`,
    {
      description: 'ICU example message',
      original_code: {
        'INTERPOLATION_1': '{{getPerson()}}',
        'INTERPOLATION_2': '{{getPlaceName()}}',
      },
      example: {
        'INTERPOLATION_1': 'Jane Doe',
        'INTERPOLATION_2': 'Paris, France',
      }
    });

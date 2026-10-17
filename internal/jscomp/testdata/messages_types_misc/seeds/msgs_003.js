const MSG_ICU_EXAMPLE = declareIcuTemplate(
    'Email: {userEmail}', // should be USER_EMAIL
    {
      description: 'Labeled email address',
      example: {
        'userEmail': 'jane@doe.com', // should be USER_EMAIL
      }
    });

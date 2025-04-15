module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [
      2,
      'always',
      [
        'feat',     // New feature
        'fix',      // Bug fix
        'docs',     // Documentation changes
        'style',    // Code style changes (formatting, etc.)
        'refactor', // Code refactoring
        'test',     // Adding or modifying tests
        'chore',    // Changes to build process or auxiliary tools
        'ci',       // CI configuration changes
        'perf',     // Performance improvements
        'revert'    // Reverting changes
      ]
    ],
    'header-max-length': [2, 'always', 100],
    'subject-case': [0],
    'type-case': [2, 'always', ['lower-case']],
    'type-empty': [2, 'never'],
    'subject-empty': [2, 'never'],
    'body-leading-blank': [1, 'always'],
    'footer-leading-blank': [1, 'always']
  },
  ignores: [
    (commit) => commit.startsWith('Merge remote-tracking branch')
  ]
}; 
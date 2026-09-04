module.exports = {
  branches: [
    { name: 'main', channel: false },
    { name: 'next', prerelease: 'rc', channel: 'next' },
  ],
  repositoryUrl: 'git@github.com:telekom/gateway-issuer-service-go.git',
  tagFormat: '${version}',
  plugins: [
    [
      '@semantic-release/commit-analyzer',
      {
        preset: 'conventionalcommits',
        releaseRules: [
          { breaking: true, release: 'major' },
          { revert: true, release: false },
        ],
      },
    ],
    'semantic-release-export-data',
    [
      '@semantic-release/release-notes-generator',
      {
        preset: 'conventionalcommits',
        presetConfig: {
          types: [
            { type: 'feat', section: 'Features', hidden: false },
            { type: 'fix', section: 'Bug Fixes', hidden: false },
            { type: 'perf', section: 'Performance Improvements', hidden: false },
          ],
        },
      },
    ],
    '@semantic-release/github',
  ],
};

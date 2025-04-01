module.exports = {
    branches: ['main'],
    repositoryUrl: 'git@github.com:telekom/gateway-issuer-service-go.git',
    plugins: [
        '@semantic-release/commit-analyzer',
        'semantic-release-export-data',
        '@semantic-release/release-notes-generator',
        '@semantic-release/changelog',
        '@semantic-release/github',
        [
            '@semantic-release/git',
            {
                assets: ['CHANGELOG.md']
            },
        ],
    ],
};
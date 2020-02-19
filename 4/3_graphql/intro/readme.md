query ($number_of_repos: Int!) {
  viewer {
    name
    login
    repositories {
      totalCount
    }
    followers {
      totalCount
    }
    starredRepositories(last: $number_of_repos) {
      totalCount
      nodes {
        name
        description
        forkCount
        homepageUrl
        stargazers {
          totalCount
        }
        updatedAt
      }
    }
  }
}



-----



query ($number_of_repos: Int!) {
  viewer {
    name
    login
    starredRepositories(last: $number_of_repos) {
      totalCount
      nodes {
        name
        description
        homepageUrl
        updatedAt
        issues(last: 3, states: OPEN) {
          edges {
            node {
              title
              url
            }
          }
        }
      }
    }
  }
}

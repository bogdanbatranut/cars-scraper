package autovit

import "carscraper/pkg/jobs"

type PostURLBuilder struct {
	criteria     jobs.Criteria
	paramsMapper ParamsMapper
}

func NewPOSTURLBuilder(criteria jobs.Criteria) *PostURLBuilder {
	return &PostURLBuilder{
		criteria:     criteria,
		paramsMapper: NewParamsMapper(),
	}
}

func (b PostURLBuilder) GetPageURL(pageNumber int) string {
	return ""
}

func (b PostURLBuilder) GetQueryPayload(pageNumber int) *AutovitPOSTPayload {
	return nil
}

//query listingScreen(
//	$after: ID,
//	$click2BuyExperimentId: String!,
//	$click2BuyExperimentVariant: String!,
//	$experiments: [Experiment!],
//	$filters: [AdvertSearchFilterInput!],
//	$includeClick2Buy: Boolean!,
//	$includePriceEvaluation: Boolean!,
//	$includePromotedAds: Boolean!,
//	$includeRatings: Boolean!,
//	$includeFiltersCounters: Boolean!,
//	$includeSortOptions: Boolean!,
//	$includeSuggestedFilters: Boolean!,
//	$itemsPerPage: Int,
//	$page: Int,
//	$parameters: [String!],
//	$searchTerms: [String!],
//	$sortBy: String,
//	$maxAge: Int,
//	$includeCepik: Boolean!)
//{\n  advertSearch
//	(\n    criteria:		{searchTerms: $searchTerms, filters: $filters}
//		\n    sortBy: $sortBy
//		\n    page: $page
//		\n    after: $after
//		\n    itemsPerPage: $itemsPerPage
//		\n    maxAge: $maxAge
//		\n    experiments: $experiments
//	\n  )
//{\n    ...advertSearchFields
//	\n    edges {
//		\n      node {
//			\n        ...lazyAdvertFields
//			\n        __typename
//			\n      }
//			\n      __typename
//		\n    }
//\n    sortOptions @include(if: $includeSortOptions) {
//	\n      searchKey
//	\n      label
//	\n      __typename
//	\n    }
//\n    __typename
//\n  }
//\n  ...Click2BuyServiceSearch
//\n  ...suggestedFilters @include(if: $includeSuggestedFilters)
//\n}
//\nfragment advertFields on Advert {
//		\n  id
//		\n  title
//		\n  createdAt
//		\n  shortDescription
//		\n  url
//		\n  badges
//		\n  category {
//			\n    id
//			\n    __typename
//			\n  }
//		\n  location {
//			\n    city {
//				\n      name
//				\n      __typename
//				\n    }
//			\n    region {
//				\n      name
//				\n      __typename
//				\n    }\n    __typename
//			\n  }
//		\n  thumbnail {
//			\n    x1: url(size: {width: 320, height: 240})
//			\n    x2: url(size: {width: 640, height: 480})
//			\n    __typename
//			\n  }
//		\n  price {
//			\n    amount {
//				\n      units
//				\n      nanos
//				\n      value
//				\n      currencyCode
//				\n      __typename
//				\n    }
//			\n    badges
//			\n    grossPrice {
//				\n      value
//				\n      currencyCode
//				\n      __typename
//				\n    }
//			\n    netPrice {
//				\n      value
//				\n      currencyCode
//				\n      __typename
//				\n    }
//			\n    __typename
//			\n  }
//		\n  parameters(keys: $parameters) {
//			\n    key
//			\n    displayValue
//			\n    label
//			\n    value
//			\n    __typename
//			\n  }
//		\n  sellerLink {
//			\n    id
//			\n    name
//			\n    websiteUrl
//			\n    logo {
//				\n      x1: url(size: {width: 140, height: 24})
//				\n      __typename
//				\n    }
//			\n    __typename
//			\n  }
//		\n  brandProgram {
//			\n    logo {
//				\n      x1: url(size: {width: 95, height: 24})
//				\n      __typename
//				\n    }
//			\n    searchUrl
//			\n    name
//			\n    __typename
//			\n  }
//		\n  dealer4thPackage {
//			\n    package {
//				\n      id
//				\n      name
//				\n      __typename
//				\n    }
//			\n    services {
//				\n      code
//				\n      label
//				\n      __typename
//				\n    }
//			\n    photos(first: 5) {
//				\n      nodes {
//					\n        url(size: {width: 644, height: 461})
//					\n        __typename
//					\n      }
//				\n      totalCount
//				\n      __typename
//				\n    }
//			\n    __typename
//			\n  }
//		\n  priceEvaluation @include(if: $includePriceEvaluation) {
//			\n    indicator
//			\n    __typename
//			\n  }
//		\n  __typename
//		\n}
//\nfragment listingAdCardFields on AdvertEdge {
//	\n  vas {
//		\n    isHighlighted
//		\n    isPromoted
//		\n    bumpDate
//		\n    __typename
//		\n  }
//	\n  node {
//		\n    ...advertFields
//		\n    __typename
//		\n  }
//	\n  __typename
//	\n}
//\nfragment advertSearchFields on AdvertSearchOutput {
//	\n  url
//	\n  sortedBy
//	\n  locationCriteriaChanged
//	\n  subscriptionKey
//	\n  totalCount
//	\n  filtersCounters @include(if: $includeFiltersCounters) {
//		\n    name
//		\n    nodes {
//			\n      name
//			\n      value
//			\n      __typename
//			\n    }
//		\n    __typename
//		\n  }
//	\n  appliedLocation {
//		\n    city {
//			\n      name
//			\n      id
//			\n      canonical
//			\n      __typename
//			\n    }
//		\n    subregion {
//			\n      name
//			\n      id
//			\n      canonical
//			\n      __typename
//			\n    }
//		\n    region {
//			\n      name
//			\n      id
//			\n      canonical
//			\n      __typename
//			\n    }
//		\n    latitude
//		\n    longitude
//		\n    mapConfiguration {
//			\n      zoom
//			\n      __typename
//			\n    }
//		\n    __typename
//		\n  }
//	\n  appliedFilters {
//		\n    name
//		\n    value
//		\n    canonical
//		\n    __typename
//		\n  }
//	\n  breadcrumbs {
//		\n    label
//		\n    url
//		\n    __typename
//		\n  }
//	\n  pageInfo {
//		\n    pageSize
//		\n    currentOffset
//		\n    __typename
//		\n  }
//	\n  facets {
//		\n    options {
//			\n      label
//			\n      url
//			\n      count
//			\n      __typename
//			\n    }
//		\n    __typename
//		\n  }
//	\n  alternativeLinks {
//		\n    name
//		\n    title
//		\n    links {
//			\n      title
//			\n      url
//			\n      counter
//			\n      __typename
//			\n    }
//		\n    __typename
//		\n  }
//	\n  latestAdId
//	\n  edges {
//		\n    ...listingAdCardFields
//		\n    __typename
//		\n  }
//	\n  topads @include(if: $includePromotedAds) {
//		\n    ...listingAdCardFields
//		\n    __typename
//		\n  }
//	\n  __typename
//	\n}
//\nfragment lazyAdvertFields on Advert {
//	\n  id
//	\n  cepikVerified @include(if: $includeCepik)
//	\n  sellerRatings(scope: PROFESSIONAL) @include(if: $includeRatings) {
//		\n    statistics {
//			\n      recommend {
//				\n        value
//				\n        suffix
//				\n        __typename
//				\n      }
//			\n      avgRating {
//				\n        value
//				\n        __typename
//				\n      }
//			\n      total {
//				\n        suffix
//				\n        value
//				\n        __typename
//				\n      }
//			\n      detailedRating {
//				\n        label
//				\n        value
//				\n        __typename
//				\n      }
//			\n      __typename
//			\n    }
//		\n    __typename
//		\n  }
//	\n  __typename
//	\n}
//\nfragment Click2BuyServiceSearch on Query {
//	\n  click2Buy @include(if: $includeClick2Buy) {
//		\n    search(
//			\n      criteria: {filters: $filters, searchTerms: $searchTerms}
//			\n      experimentId: $click2BuyExperimentId
//			\n      experimentVariant: $click2BuyExperimentVariant
//			\n      itemsPerPage: 4
//			\n      page: $page\n    ) {
//			\n      __typename
//			\n      ... on Click2BuySearchOutput {
//				\n        edges {
//					\n          node {
//						\n            id
//						\n            title
//						\n            shortDescription
//						\n            url
//						\n            detailUrl
//						\n            sellerLink {
//							\n              id
//							\n              name
//							\n              websiteUrl
//							\n              __typename
//							\n            }
//						\n            photos {
//							\n              url
//							\n              __typename
//							\n            }
//						\n            chips
//						\n            checks
//						\n            utmContentType
//						\n            carDetails {
//							\n              year
//							\n              mileage
//							\n              engineVolume
//							\n              fuelType
//							\n              make
//							\n              model
//							\n              __typename
//							\n            }
//						\n            price {
//							\n              amount {
//								\n                value
//								\n                currencyCode
//								\n                __typename
//								\n              }
//							\n              __typename
//							\n            }
//						\n            installmentPrice {
//							\n              amount {
//								\n                value
//								\n                currencyCode
//								\n                __typename
//								\n              }
//							\n              __typename
//							\n            }
//						\n            __typename
//						\n          }
//					\n          __typename
//					\n        }
//				\n        __typename
//				\n      }
//			\n    }
//		\n    __typename
//		\n  }
//	\n  __typename
//	\n}
//\nfragment suggestedFilters on Query {
//	\n  suggestedFilters(criteria: {searchTerms: $searchTerms, filters: $filters}) {
//		\n    key
//		\n    name
//		\n    values {
//			\n      value
//			\n      appliedFilters {
//				\n        name
//				\n        value
//				\n        __typename
//				\n      }
//			\n      __typename
//			\n    }
//		\n    __typename
//		\n  }
//	\n  __typename
//	\n}

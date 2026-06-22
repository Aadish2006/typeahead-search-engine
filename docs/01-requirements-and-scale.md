# Problem Statement
> Design a typeahead application.

A typeahead is a specific form of an autocomplete which is meant to type ahead of the user, thereby
saving the user's time of typing the rest of the query. It is not meant to suggest new things, rather
predict the query that user wants to type, hence the name typeahead.

The list of suggestions are the most popular queries made by users of the search system for which
the typeahead system is supposed to be built on.

# Functional Requirements
1. As the user is typing in the search-box, we should show type ahead suggestions.
2. The suggestions shown are the searches that other people made in the past.
3. Limit the suggestions to show only top 10, relevant suggestions.
    - _relevant_: any suggestion is relevant if the prefix of the suggestion matches with the query.
    - _top 10_: we need to keep track the number of times a search query has been searched for globally. Given a partial query we will show the search queries which match the prefix with partial query and have the highest frequency.
4. With every character that the user types, we need to show new suggestions based on the newly typed query - send a new API request on character press.
5. We will start showing typeahead suggestions only after the user has typed atleast 3 characters in the search bar.

# Non Functional Requirements
> The only non functional requirement is that the typeahead suggestions should be optimized for **low latency**.

### What does eventual consistency mean in this context?
It means that if a user makes a search query then it will lazily/asynchronously update the frequency count for that specific query.
With eventual consistency a user will not see the absolute updated data, because it won't account for the searches made by other users
which have not yet been added to the database i.e. the user will perform **stale reads** when seeing suggestions.

### Can we afford stale reads?
Definitely. Why? because the entire typeahead system is not the core feature. The core feature is the search engine which is already implemented, typeahead is just a good to have feature on top of the core (search engine) system. Stale reads are not a problem because we are not displaying the
frequency count of the query, we only need to the relative order of the top 10 relevant queries.

### What does data loss mean in this context?
It means that if a user makes a search query, then the database does not update the frequency for that specific query. In case of data loss the frequencies for some queries will be off by some count.

### Can we afford data loss?
Most definitely. Why? because as we discussed previously the entire typeahead system is a good to have and is not necessarily a must have feature. Also if there is some amount of data loss that is absolutely fine, since we only need the trends for the queries and not the absolute precise frequency to compare the suggestions.

### Consistency or Latency?
**Latency**. Why? As discussed above, we can sacrifice consistency to get a low latency system.

### During a network partition: Consistency or Availability?
**Availability**. Same reasons as discussed above.

# Scale Estimation

## Requests Estimation
We are going to assume that the typeahead system is being made for Google's search engine.

The population of the entire world is 8 billion, where roughly 5 billion people use the Internet. It is safe to assume that every user of Internet is also a user of Google, due to its sheer dominance and popularity.

**Daily Active Users (DAU)**: 5 billion

By using the Paretto's principle we can estimate that 20% of the users are going to be active users.

**Daily Actually Active Users**: 20% of 5 billion = 1 billion <br/>
**Search queries per active-user per day**: 20 <br/>
**Total searches per day**: 20 * 1 billion = 20 billion searches / day <br/>
**Searches per second**: 20 billion / 10^5 = **200,000 searches / second**.

Google search engine processes 200,000 searches per second, the typeahead system functions before the query is performed, when the user is typing the query. We need to show suggestions on every character press, which means we will have to make a request on every character press. So the typeahead system will have way more requests per second than the actual search engine.

**Average number of characters in a search query**: 10 characters

We start showing the suggestions after typing 3 characters, which means the number of requests to the typeahead system is 3 less than the average number of characters in a search query i.e. 10 - 3 = 7 which is roughly equal to 10.

**Number of requests to typeahead per query**: 10 <br/>
**Total requests to typeahead per second**: 10 * 200,000 = **2 Million requests / second**

## Data Estimation

### Schema

| Search Query | Frequency |
| :-- | :-- |
| Why is Suryakumar Yadav in the playing 11? | 10782 |
| What are LSM trees? | 7411 |
| What is the capital of India? | 3498 |

**Size of one row**: 10 bytes (average length of search query is 10) + 8 bytes (for frequency) ~ 20 bytes
**Number of new rows per day**: number of unique search queries

**Percentage of unique search queries**: 10%
**Unique search queries per day**: 0.1 * 20 billion searches / day = 2 billion unique searches / day
**Total data generated per day**: 2 billion unique searches / day * 20 bytes per row = **40 GB / day**

**Total data for a period of 20 years**: 40 GB / day * 20 years * 400 days / year = **320 TB**

### Can we store all the data in a single server?
Maybe, but one single server will definitely not be enough to handle 2 million requests / second. This means we need to shard the data into multiple database servers.

### Read and Write Metrics

Reads:  2 Million requests / second
Writes: 200,000 requests / second

The system is both read and write heavy. It is not theoretically possible to create a system which offers low latency for both a read and write heavy workload. The reason behind it is that there is no database system which is optimized for both, we can only optimize it for any one.

> We can engineer a solution where we absorb a large portion of reads from a cache and optimize the database for writes.

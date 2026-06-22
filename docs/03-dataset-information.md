## Dataset Information

This assignment uses the [AOL Search Logs Dataset](https://archive.org/download/academictorrents_cd339bddeae7126bb3b15f3a72c903cb0c401bd1/AOL_search_data_leak_2006.zip). It is a collection of real searches made by people of America between March 1, 2006 to May 31, 2006. During this period of 3 months, there were about **20 million** search queries made by approximately **650,000** anonymized users. Since the data comes from the real world, it is highly practical for a type ahead system.

The dataset contains many columns of search data, but for this assignment we only need the _query_ column which is the actual user query that people made on the search engines. The frequency needs to be calculated for each unique query and create a specialized dataset for the use case of making a performant typeahead system.

Since the queries are typed by people, naturally there will be many typos. For example if many people are trying to search for porn.com then many queries would have typos like pon.com or por.com. I am not normalizing these types of queries since this features is not part of the MVP, but can be included in future requirements.

Fun Fact: The dataset was released by AOL Research for academic purposes in 2006, but ended up causing a massive scandal. Even thought the user ID for each user were anonymized, journalists and researchers were still able to cross-reference the search queries to specific individuals.

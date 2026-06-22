const suggestionList = document.getElementById("suggestion-list");
const searchBarInput = document.getElementById("search-bar-input");

let activeRequestController = null;

searchBarInput.addEventListener("input", loadSuggestions);
searchBarInput.addEventListener("keydown", handleSearchInput)
searchBarInput.value = "";
suggestionList.replaceChildren();

async function loadSuggestions(inputEvent) {
  // Retrieve partial search query from user input.
  const searchPrefix = inputEvent.target.value;

  // Clear suggestions if input is empty or exceeds max length.
  if (searchPrefix === "" || searchPrefix.length > 50) {
    suggestionList.classList.remove("visible");
    suggestionList.replaceChildren();
    return;
  }

  // Cancel any in-flight requests.
  if (activeRequestController) {
    activeRequestController.abort();
  }
  activeRequestController = new AbortController();

  try {
    // Fetch suggestions from the backend API.
    const suggestions = await fetchSuggestions(searchPrefix, 10, activeRequestController.signal);

    // Clear previous suggestions from the list.
    suggestionList.replaceChildren();

    // Populate the suggestion list with results.
    for (const suggestion of suggestions) {
      if (suggestion.length < 50) {
        const listItem = document.createElement('li');
        listItem.innerHTML = suggestion.substring(0, searchPrefix.length) + "<b>" + suggestion.substring(searchPrefix.length) + "</b>";
        suggestionList.appendChild(listItem);
      }
    }
    suggestionList.classList.add("visible");
  } catch (error) {
    if (error.name !== "AbortError") {
      console.error("Failed to fetch suggestions:", error);
    }
  }
}

async function fetchSuggestions(prefix, limit, signal) {
  const url = new URL("http://localhost:8080/api/v1/suggestions")
  url.searchParams.append("prefix", prefix)
  url.searchParams.append("limit", limit)

  const response = await fetch(url, { signal });
  const suggestions = await response.json();
  return suggestions;
}

async function handleSearchInput(event) {
  if (event.key !== "Enter") {
    return;
  }

  const query = search_bar.value;
  if (query === "") {
    return
  }

  await fetch("http://127.0.0.1:8080/api/v1/search", {
    method: "POST",
    body: query
  });

  search_bar.value = '';
  suggestion_list.replaceChildren();
}

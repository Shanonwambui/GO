```
func deleteMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID  == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
}
```
movies[:index]: This creates a new slice containing all elements from the start of the movies slice up to, but not including, the element at index.

movies[index+1:]: This creates another slice containing all elements from the movies slice starting just after the index element to the end of the slice.

append(movies[:index], movies[index+1:]...): The append function concatenates the two slices created in steps 1 and 2. The ... operator is used to unpack the elements of the second slice so that they can be appended individually.

movies = append(...): The result of the append operation is assigned back to the movies slice, effectively removing the element at the specified index.
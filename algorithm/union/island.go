package union

// leetcode 200
func Island(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	var res int
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 {
				res++
				dfs(grid, i, j, m, n)
			}
		}
	} // m 行
	return res
}

func dfs(gird [][]int, i, j, m, n int) {
	if i < 0 || i >= m || j < 0 || j >= n || gird[i][j] != 1 {
		return
	}
	gird[i][j] = 2
	dfs(gird, i+1, j, m, n)
	dfs(gird, i-1, j, m, n)
	dfs(gird, i, j+1, m, n)
	dfs(gird, i, j-1, m, n)
}

<script setup>
import { ref, onMounted } from 'vue'
import { CreateUser, GetAllUsers, GetUser, UpdateUser, DeleteUser, SearchUsers } from '../wailsjs/go/main/App.js'

// 響應式數據
const users = ref([])
const loading = ref(false)
const message = ref('')
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 表單數據
const formData = ref({
  name: '',
  email: '',
  age: 0
})

const editingUser = ref(null)

// 載入所有用戶
const loadUsers = async () => {
  loading.value = true
  try {
    const result = await SearchUsers(searchKeyword.value, currentPage.value, pageSize.value)
    users.value = result.users
    total.value = result.total
  } catch (error) {
    message.value = `載入用戶失敗: ${error}`
  } finally {
    loading.value = false
  }
}

// 創建用戶
const createUser = async () => {
  if (!formData.value.name || !formData.value.email) {
    message.value = '請填寫姓名和電子郵件'
    return
  }
  
  try {
    await CreateUser(formData.value.name, formData.value.email, formData.value.age)
    message.value = '用戶創建成功'
    formData.value = { name: '', email: '', age: 0 }
    loadUsers()
  } catch (error) {
    message.value = `創建用戶失敗: ${error}`
  }
}

// 更新用戶
const updateUser = async () => {
  if (!editingUser.value) return
  
  try {
    await UpdateUser(editingUser.value.id, formData.value.name, formData.value.email, formData.value.age)
    message.value = '用戶更新成功'
    editingUser.value = null
    formData.value = { name: '', email: '', age: 0 }
    loadUsers()
  } catch (error) {
    message.value = `更新用戶失敗: ${error}`
  }
}

// 刪除用戶
const deleteUser = async (id) => {
  if (!confirm('確定要刪除這個用戶嗎？')) return
  
  try {
    await DeleteUser(id)
    message.value = '用戶刪除成功'
    loadUsers()
  } catch (error) {
    message.value = `刪除用戶失敗: ${error}`
  }
}

// 編輯用戶
const editUser = (user) => {
  editingUser.value = user
  formData.value = {
    name: user.name,
    email: user.email,
    age: user.age
  }
}

// 取消編輯
const cancelEdit = () => {
  editingUser.value = null
  formData.value = { name: '', email: '', age: 0 }
}

// 搜尋
const search = () => {
  currentPage.value = 1
  loadUsers()
}

// 分頁
const changePage = (page) => {
  currentPage.value = page
  loadUsers()
}

// 組件掛載時載入數據
onMounted(() => {
  loadUsers()
})
</script>

<template>
  <div class="container">
    <h1>MongoDB 用戶管理系統</h1>
    
    <!-- 訊息顯示 -->
    <div v-if="message" class="message" :class="{ error: message.includes('失敗') }">
      {{ message }}
    </div>
    
    <!-- 表單 -->
    <div class="form-section">
      <h2>{{ editingUser ? '編輯用戶' : '新增用戶' }}</h2>
      <form @submit.prevent="editingUser ? updateUser() : createUser()">
        <div class="form-group">
          <label>姓名:</label>
          <input v-model="formData.name" type="text" required />
        </div>
        <div class="form-group">
          <label>電子郵件:</label>
          <input v-model="formData.email" type="email" required />
        </div>
        <div class="form-group">
          <label>年齡:</label>
          <input v-model.number="formData.age" type="number" min="0" />
        </div>
        <div class="form-actions">
          <button type="submit" :disabled="loading">
            {{ editingUser ? '更新' : '創建' }}
          </button>
          <button v-if="editingUser" type="button" @click="cancelEdit">取消</button>
        </div>
      </form>
    </div>
    
    <!-- 搜尋 -->
    <div class="search-section">
      <div class="search-group">
        <input v-model="searchKeyword" type="text" placeholder="搜尋用戶..." />
        <button @click="search">搜尋</button>
        <button @click="searchKeyword = ''; search()">清除</button>
      </div>
    </div>
    
    <!-- 用戶列表 -->
    <div class="users-section">
      <h2>用戶列表</h2>
      <div v-if="loading" class="loading">載入中...</div>
      <div v-else-if="!users || users.length === 0" class="no-data">沒有找到用戶</div>
      <table v-else class="users-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>姓名</th>
            <th>電子郵件</th>
            <th>年齡</th>
            <th>創建時間</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.name }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.age }}</td>
            <td>{{ new Date(user.created_at).toLocaleString() }}</td>
            <td>
              <button @click="editUser(user)" class="btn-edit">編輯</button>
              <button @click="deleteUser(user.id)" class="btn-delete">刪除</button>
            </td>
          </tr>
        </tbody>
      </table>
      
      <!-- 分頁 -->
      <div v-if="total > pageSize" class="pagination">
        <button 
          v-for="page in Math.ceil(total / pageSize)" 
          :key="page"
          @click="changePage(page)"
          :class="{ active: page === currentPage }"
        >
          {{ page }}
        </button>
      </div>
    </div>
  </div>
</template>

<style>
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  font-family: Arial, sans-serif;
}

h1, h2 {
  color: #333;
  margin-bottom: 20px;
}

.message {
  padding: 10px;
  margin: 10px 0;
  border-radius: 4px;
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.message.error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.form-section, .search-section, .users-section {
  background: #f8f9fa;
  padding: 20px;
  margin: 20px 0;
  border-radius: 8px;
  border: 1px solid #dee2e6;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: black;
}

.form-group input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
}

.form-actions {
  display: flex;
  gap: 10px;
}

.search-group {
  display: flex;
  gap: 10px;
  align-items: center;
}

.search-group input {
  flex: 1;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

button:not(:disabled):hover {
  opacity: 0.8;
}

.btn-edit {
  background-color: #007bff;
  color: white;
  margin-right: 5px;
}

.btn-delete {
  background-color: #dc3545;
  color: white;
}

.loading, .no-data {
  text-align: center;
  padding: 20px;
  color: #666;
}

.users-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 10px;
  color: black;
}

.users-table th,
.users-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

.users-table th {
  background-color: #f8f9fa;
  font-weight: bold;
}

.users-table tr:hover {
  background-color: #f5f5f5;
}

.pagination {
  display: flex;
  justify-content: center;
  gap: 5px;
  margin-top: 20px;
}

.pagination button {
  padding: 5px 10px;
  border: 1px solid #ddd;
  background-color: white;
}

.pagination button.active {
  background-color: #007bff;
  color: white;
  border-color: #007bff;
}
</style>

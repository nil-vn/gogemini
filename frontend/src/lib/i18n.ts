import { writable, derived, get } from 'svelte/store';

export type Locale = 'en' | 'vi';

type Dict = Record<string, string>;

const dictionaries: Record<Locale, Dict> = {
  en: {
    appTitle: 'GoGemini Admin',
    navDashboard: 'Dashboard',
    navSettings: 'Settings',
    navLogout: 'Logout',
    dashboardTitle: 'Dashboard metrics',
    managementTitle: '{module} management',
    filterLabel: 'Filter',
    filterPlaceholder: 'Filter',
    sortById: 'Sort by id',
    create: 'Create',
    update: 'Update',
    reset: 'Reset',
    prev: 'Prev',
    next: 'Next',
    page: 'Page',
    login: 'Login',
    admin: 'Admin',
    recordUpdated: 'Record updated',
    recordCreated: 'Record created',
    recordDeleted: 'Record deleted',
    deleteConfirm: 'Delete {module} #{id}?',
    tableList: '{title} list',
    tableTotal: 'Total: {total} | page: {page} | page size: {pageSize}',
    noRecords: 'No records',
    id: 'id',
    data: 'data',
    actions: 'actions',
    detail: 'Detail',
    delete: 'Delete',
    settingsTitle: 'System settings',
    currency: 'Currency',
    theme: 'Theme',
    language: 'Language',
    save: 'Save',
    searchLabel: 'Search',
    searchPlaceholder: 'Search all modules',
    searchResults: 'Global search results',
    uploadImage: 'Upload image',
    upload: 'Upload',
    uploadSuccess: 'Image uploaded successfully',
    dark: 'dark',
    light: 'light'
  },
  vi: {
    appTitle: 'GoGemini Admin',
    navDashboard: 'Tổng quan',
    navSettings: 'Thiết lập',
    navLogout: 'Đăng xuất',
    dashboardTitle: 'Chỉ số dashboard',
    managementTitle: 'Quản lý {module}',
    filterLabel: 'Lọc',
    filterPlaceholder: 'Lọc',
    sortById: 'Sắp xếp theo id',
    create: 'Tạo mới',
    update: 'Cập nhật',
    reset: 'Đặt lại',
    prev: 'Trước',
    next: 'Sau',
    page: 'Trang',
    login: 'Đăng nhập',
    admin: 'Quản trị',
    recordUpdated: 'Đã cập nhật bản ghi',
    recordCreated: 'Đã tạo bản ghi',
    recordDeleted: 'Đã xóa bản ghi',
    deleteConfirm: 'Xóa {module} #{id}?',
    tableList: 'Danh sách {title}',
    tableTotal: 'Tổng: {total} | trang: {page} | kích thước trang: {pageSize}',
    noRecords: 'Không có bản ghi',
    id: 'id',
    data: 'dữ liệu',
    actions: 'thao tác',
    detail: 'Chi tiết',
    delete: 'Xóa',
    settingsTitle: 'Thiết lập hệ thống',
    currency: 'Tiền tệ',
    theme: 'Giao diện',
    language: 'Ngôn ngữ',
    save: 'Lưu',
    searchLabel: 'Tìm kiếm',
    searchPlaceholder: 'Tìm toàn bộ module',
    searchResults: 'Kết quả tìm kiếm tổng',
    uploadImage: 'Tải ảnh lên',
    upload: 'Tải lên',
    uploadSuccess: 'Tải ảnh thành công',
    dark: 'tối',
    light: 'sáng'
  }
};

export const locale = writable<Locale>('en');
export const t = derived(locale, ($locale) => (key: string, vars: Record<string, string | number> = {}) => {
  const dict = dictionaries[$locale] ?? dictionaries.en;
  const template = dict[key] ?? dictionaries.en[key] ?? key;
  return Object.entries(vars).reduce((acc, [k, v]) => acc.split(`{${k}}`).join(String(v)), template);
});

export function setLocale(next: string) {
  if (next === 'vi' || next === 'en') locale.set(next);
}

export function getTranslator() { return get(t); }

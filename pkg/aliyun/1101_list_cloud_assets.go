package aliyun

// list_cloud_assets

func ListCloudAssets() {
	OSSListBuckets()
	ECSListInstances()
	RAMListUsers()
	RAMListRoles()
	DomainListDomains()
}

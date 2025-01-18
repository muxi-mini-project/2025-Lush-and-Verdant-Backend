package dao

import "2025-Lush-and-Verdant-Backend/model"

// 对slogan库的初始化
func CreateSlogan() {
	db := NewDB(dsn)

	db.AutoMigrate(&model.Slogan{})

	// 对slogan库的初始化
	slogans := []model.Slogan{
		{ID: 1, Slogan: "没有什么做不到，只有你想不到"},
		{ID: 2, Slogan: "走自己的路，活出自己的精彩人生"},
		{ID: 3, Slogan: "成功一定有方法，失败一定有原因"},
		{ID: 4, Slogan: "有志者自有千方百计，无志者只有千难万难"},
		{ID: 5, Slogan: "生如蝼蚁当立鸿鹄之志，命如薄纸应有不屈之心"},
		{ID: 6, Slogan: "少年何妨梦摘星，敢挽桑弓射玉衡"},
		{ID: 7, Slogan: "宁如飞萤赴火，不做樗木长春"},
		{ID: 8, Slogan: "风雪压我两三年，我笑风轻雪如棉"},
		{ID: 9, Slogan: "海到无边天作岸，山登绝顶我为峰"},
		{ID: 10, Slogan: "好习惯是这一生享不尽的财富，坏习惯是一辈子还不清的债务"},
	}

	db.Create(&slogans)

}

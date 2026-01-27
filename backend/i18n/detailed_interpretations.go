package i18n

import "star/models"

// GetDetailedInterpretation returns detailed interpretation for an event
func (t *Translator) GetDetailedInterpretation(eventType string, planet1, planet2 models.PlanetID, aspect, house string, isPositive bool) string {
	// For transit house events
	if eventType == "transit_house" {
		return t.getTransitHouseInterpretation(planet1, house)
	}
	
	// For progression events
	if eventType == "secondary_progression" || eventType == "tertiary_progression" {
		return t.getProgressionInterpretation(planet1, planet2, aspect, isPositive)
	}
	
	// For aspect events
	if eventType == "aspect" {
		return t.getAspectInterpretation(planet1, planet2, aspect, isPositive)
	}
	
	return ""
}

// getTransitHouseInterpretation returns detailed interpretation for transit houses
func (t *Translator) getTransitHouseInterpretation(planet models.PlanetID, house string) string {
	key := string(planet) + "_house_" + house
	
	switch t.lang {
	case Chinese:
		return getChineseTransitHouseInterpretation(key, house)
	case Russian:
		return getRussianTransitHouseInterpretation(key, house)
	default:
		return getEnglishTransitHouseInterpretation(key, house)
	}
}

// getProgressionInterpretation returns detailed interpretation for progressions
func (t *Translator) getProgressionInterpretation(p1, p2 models.PlanetID, aspect string, isPositive bool) string {
	key := string(p1) + "_" + aspect + "_" + string(p2)
	
	switch t.lang {
	case Chinese:
		return getChineseProgressionInterpretation(key, isPositive)
	case Russian:
		return getRussianProgressionInterpretation(key, isPositive)
	default:
		return getEnglishProgressionInterpretation(key, isPositive)
	}
}

// getAspectInterpretation returns detailed interpretation for aspects
func (t *Translator) getAspectInterpretation(p1, p2 models.PlanetID, aspect string, isPositive bool) string {
	key := string(p1) + "_" + aspect + "_" + string(p2)
	
	switch t.lang {
	case Chinese:
		return getChineseAspectInterpretation(key, isPositive)
	case Russian:
		return getRussianAspectInterpretation(key, isPositive)
	default:
		return getEnglishAspectInterpretation(key, isPositive)
	}
}

// ========== Chinese Transit House Interpretations ==========

func getChineseTransitHouseInterpretation(key string, house string) string {
	interpretations := map[string]string{
		// Sun through houses - 重新设计以强调维度
		
		// 1宫：健康+事业（双维度）
		"sun_house_1": "这是个人魅力和生命活力全面提升的时期。你会感到**身体充满能量**，精神状态极佳，**体力和耐力都处于高峰**。这股强大的生命力让你在**职业场合表现出色**，**领导能力和个人魅力**吸引他人注意。适合在工作中主动出击，展现你的**专业能力和健康活力**。注意保持良好的**作息习惯**，让充沛的精力为**职业发展**助力。",
		
		// 2宫：财运（单维度）
		"sun_house_2": "**财务和物质安全**成为关注的焦点。这段时间你会更加重视**收入来源**和**资产状况**，可能会有**增加收入**的机会或是对**财务规划**进行重新评估。你对自己的价值和能力有了更清晰的认识，这有助于提升**赚钱能力**。适合**制定理财计划、寻找新的收入来源**，或是**投资于能够长期增值**的事物。审视个人**金钱观**，确立什么样的**物质基础**对你真正重要。",
		
		// 3宫：事业+关系（双维度）
		"sun_house_3": "沟通表达为你带来**职业机遇**和**人际拓展**。你会发现**工作中的交流合作**变得频繁，**业务洽谈、会议展示**的机会增多。同时**社交网络**也在扩大，**结识新朋友、维系旧关系**都很顺利。这是在**职场建立人脉**、通过**沟通技巧推进项目**的好时机。短途出差可能带来**商业合作**，**社交活动**也能助力**事业发展**。注意在**工作沟通**和**私人社交**间保持平衡。",
		
		// 4宫：关系（单维度）
		"sun_house_4": "**家庭关系和情感纽带**成为生活的重心。你会更多地关注**与家人的相处**，渴望建立或改善居住环境，**亲密关系的质量**变得更加重要。这是**修复家庭关系、增进家人感情**的好时期。你可能会回顾过去，处理**童年家庭模式**对现在**亲密关系**的影响。适合花更多时间**陪伴家人、深化情感连接**，或是**改善居住环境让家更温馨**。内心对**归属感和情感支持**的需求需要得到满足。",
		
		// 5宫：关系+灵性（双维度）
		"sun_house_5": "**爱情创造力**和**内在表达**成为生活的亮点。这段时间你在**浪漫关系**中更有魅力，**恋爱运势**上升，容易开始新恋情或**深化现有感情**。同时你的**创造天赋**被激活，通过**艺术、娱乐、游戏**来表达**真实自我**。这是追求**浪漫爱情**、从事**创意工作**、或是参与**自我表达类活动**的绝佳时期。你的个人魅力和**真实性**吸引他人。适合通过**创造性活动探索内心**，在**亲密关系中展现真我**。",
		
		// 6宫：健康+事业（双维度）
		"sun_house_6": "**健康管理**和**工作效率**成为关注重点。这是建立**良好生活习惯**、改善**工作流程**的好时机。你会更加注重**身体健康**，可能开始新的**健身计划**或调整**饮食作息**。同时**工作中的表现**提升，**职业技能**得到发展，你可能承担更多**工作职责**或是改进**业务方法**。适合**体检养生、学习新技能、优化工作安排**。**健康的身体**支撑**高效的工作**，**规律的作息**带来**职业上的稳定表现**。",
		
		// 7宫：关系（单维度）
		"sun_house_7": "**人际关系和伙伴合作**成为焦点。这段时间你会更多地通过**与他人的互动**来认识自己，重要的**一对一关系**得到发展。可能会遇到重要的**合作伙伴**或是**深化现有的亲密关系**。这是**签订合作协议、建立伙伴关系**、或是**改善婚姻感情**的好时期。你会学习**妥协、平衡和公平对待他人**。适合**寻求情感咨询、解决关系冲突**，或是开始新的**合作关系**。注意在**关系中保持自我**，不要过度依赖对方。",
		
		// 8宫：灵性+财运（双维度）
		"sun_house_8": "**深层心理转化**和**共享资源管理**成为主题。这段时间你会面对**深层次的内在探索**，可能经历**灵性上的蜕变和重生**。同时**财务方面涉及他人的资源**，如**贷款、投资、遗产、共同财产**等需要处理。**亲密关系的深度**和**金钱的共享**都在考验你。这是进行**心理疗愈、冥想修行**，同时**处理债务、整理财务**的好时期。**面对内心阴影**会带来**灵性成长**，**妥善管理共同财产**带来**物质安全**。",
		
		// 9宫：灵性（单维度）
		"sun_house_9": "**精神视野开阔**和**智慧探索**成为追求。这段时间你渴望扩展**哲学认知、宗教信仰、人生意义**等**精神层面**的理解。可能会有长途**灵性之旅**的机会，或是开始**高等教育、深度学习**。你对**生命意义和终极真理**的思考增多，寻求更高层次的**精神觉醒**。这是**学习冥想、研读哲学、探索信仰**的好时期。适合规划**精神成长路径**、探索不同的**灵性实践方式**。保持开放的心态，接纳新的**精神体验**和**智慧启迪**。",
		
		// 10宫：事业（单维度）
		"sun_house_10": "**事业发展**和**职业成就**成为关注的中心。这段时间你的**职业生涯**和**社会地位**受到重视，可能会有**晋升机会**或是承担更多**工作责任**。你的**专业能力**和**业绩成果**得到公众认可，适合**设定职业目标、推进重要项目**。与**上级领导**或**行业权威**的关系变得重要。这是建立**职业声誉**、实现**事业抱负**的关键时期。适合**展现专业实力**，在**工作中追求卓越**，但也要注意保持谦逊和对团队的尊重。",
		
		// 11宫：关系+灵性（双维度）
		"sun_house_11": "**社群友谊**和**集体理想**成为焦点。这段时间你会更多地参与**团体活动**，结交**志同道合的朋友**，**社交圈子**扩大。同时你对**社会愿景**和**人类理想**的关注增加，渴望为更大的**精神目标**做出贡献。这是**发展有意义的友谊、加入有共同信念的团体**的好时期。**朋友圈**带来**精神启发**，**集体活动**促进**灵性成长**。适合通过**社交网络探索理想**，在**志同道合的圈子**中**提升意识层次**。",
		
		// 12宫：灵性（单维度）
		"sun_house_12": "**内在灵性探索**和**潜意识疗愈**成为主题。这段时间你会更多地**独处静心**，进行**冥想反思**和**精神修行**。**潜意识的内容**浮现，可能通过**梦境、直觉、灵感**获得**精神启示**。这是**疗愈旧伤、释放业力、准备灵性重生**的时期。适合进行**冥想打坐、瑜伽修行、能量疗愈**，或是**处理隐藏的心理议题**。你可能会感到需要**从外界退隐**，这是**灵魂的自然需求**。**面对内心阴影**会带来**深层的灵性转化**。",
		
		// Moon through houses
		"moon_house_1": "这两天你可能感到较强的情感波动。可能会变得喜形于色，外在行为和表现都受情绪影响。你会发现自己对外界的反应更加敏感，也更加容易受到他人的情绪感染。这段时间适合表达真实的感受，但也要注意情绪管理，避免过度反应。你的直觉力增强，可以依靠内心的感觉做决定。",
		
		"moon_house_2": "情感安全感与物质安全感紧密相连。你会更加重视稳定和舒适，可能会通过购物或美食来满足情感需求。这段时间适合整理财务、评估资源，或是投资于能带来安全感的事物。你对美的事物特别敏感，可以享受感官的愉悦。注意不要过度消费来填补情感空虚。",
		
		"moon_house_3": "情绪表达和沟通变得活跃。你会有更多的交流需求，想要分享自己的感受和想法。与兄弟姐妹、邻居的互动增多，短途出行频繁。这段时间适合写日记、阅读、学习新知识。你的思维受情感影响，可能更容易感性地看待问题。保持沟通的真诚和同理心。",
		
		"moon_house_4": "家的归属和情感根基成为关注点。你会渴望待在熟悉的环境中，与家人的情感连接加深。童年记忆和家庭模式可能浮现，影响当前的情绪。这段时间适合照顾家人、整理家居、或是处理家庭关系。你需要情感上的滋养和安全感。",
		
		"moon_house_5": "情感表达变得热情和创造性。你渴望浪漫、娱乐和创造性的活动。与孩子的互动或是个人爱好能带来情感满足。这段时间适合追求快乐、表达爱意、从事艺术创作。你的情绪可能比较戏剧化，但也充满活力和热情。",
		
		"moon_house_6": "情绪影响身体健康和日常工作。你会更加关注健康习惯，可能通过规律的作息来稳定情绪。工作中需要情感支持，与同事的关系影响工作状态。这段时间适合调整生活习惯、关注健康、建立日常规律。服务他人能带来情感满足。",
		
		"moon_house_7": "情感需求通过一对一关系得到满足。你会更加依赖亲密伴侣，渴望情感上的共鸣和理解。人际关系中的情感动态变得明显，你可能更容易受到他人情绪的影响。这段时间适合深化亲密关系、寻求伙伴的支持，但也要保持情感独立。",
		
		"moon_house_8": "情绪进入深层，你会面对强烈的情感体验。可能会经历情感的蜕变或是深层的心理净化。与他人的情感纽带加深，共享资源和亲密度增加。这段时间适合探索内心深处、面对恐惧、进行情感疗愈。情绪可能比较强烈，需要安全的表达方式。",
		
		"moon_house_9": "情感扩展到更广阔的领域。你渴望通过学习、旅行来满足情感需求，对不同文化和哲学产生兴趣。情绪变得乐观和开放，但也可能不太稳定。这段时间适合探索新的信念、拓展视野，通过精神成长来滋养情感。",
		
		"moon_house_10": "公众形象和情感需求产生交集。你的情绪可能在公共场合表现出来，职业生涯受到情感的影响。你渴望在事业上获得认可来满足情感需求。这段时间适合在职业中展现真实的自己，但也要注意职业形象的管理。",
		
		"moon_house_11": "友谊和社群活动满足情感需求。你会更多地参与团体活动，从朋友那里获得情感支持。对未来的愿景带来情感上的激励。这段时间适合发展友谊、参与社群、追求集体目标。情感上需要归属感和共同理想。",
		
		"moon_house_12": "情绪进入潜意识层面，你可能会感到情感上的模糊和敏感。需要独处来处理内在的情绪，梦境和直觉变得活跃。这段时间适合冥想、独处、进行情感疗愈。你可能会感到情绪脆弱，需要给自己空间和时间来恢复。",
		
		// Venus through houses - key ones
		"venus_house_5": "浪漫和创造力达到高峰。这段时间充满魅力和吸引力，爱情关系可能开始或升温。你会更加享受生活的乐趣，艺术创作和娱乐活动带来愉悦。与孩子的互动也充满欢乐。这是追求浪漫、表达爱意、享受创造性活动的美好时期。",
		
		"venus_house_7": "关系和谐成为焦点。你会更加重视伙伴关系，渴望建立平衡和美好的连接。这是改善亲密关系、建立合作伙伴关系的好时机。你的外交技巧和魅力帮助你在人际关系中如鱼得水。适合解决关系冲突、签订合约。",
		
		"venus_house_11": "社交魅力大绽放，友谊和社群活动带来愉悦。你会结交新朋友，参与有趣的团体活动。社交网络扩大，可能通过朋友圈遇到浪漫对象。这段时间适合参与社交活动、追求集体目标、享受友谊的美好。",
		
		// Mars through houses - key ones
		"mars_house_1": "行动力和个人主动性大大增强。你会感到精力充沛，想要开始新的项目或挑战。自信心提升，但也可能变得冲动。这是实现个人目标、展现领导力的好时机，但要注意控制脾气，避免与他人冲突。",
		
		"mars_house_10": "事业野心和职业推动力达到高峰。你会全力以赴追求职业目标，可能会有重要的职业突破。工作中展现出强大的执行力和竞争力。这是争取晋升、推进重要项目的好时机，但要注意与权威人物的关系。",
		
		"mars_house_11": "积极勇敢地追求理想和社会目标。你会参与团体活动，为共同的目标而奋斗。社交中展现出行动力和领导能力。这段时间适合推动社会变革、实现集体目标，但要注意团队合作，避免独断专行。",
		
		// Jupiter through houses - key ones
		"jupiter_house_2": "财富与机遇增加的时期。收入可能提升，或是获得新的赚钱机会。你对物质的态度变得更加乐观和慷慨。这是投资、扩展财务资源的好时机。注意不要过度消费或过于乐观，保持理性的财务规划。",
		
		"jupiter_house_9": "智慧启蒙和精神成长的黄金时期。你会对哲学、宗教、高等教育产生浓厚兴趣。可能有机会远行或接触不同文化。视野开阔，对生命意义有更深的理解。这是学习、旅行、追求真理的绝佳时机。",
		
		"jupiter_house_11": "贵人现身还请抓住机会。友谊和社交活动带来幸运和机遇。你会结识有影响力的朋友，社交网络扩大。参与团体活动能带来成长和收获。这是实现长期目标、获得社会支持的好时期。",
		
		// Saturn through houses - key ones
		"saturn_house_1": "自律和责任感增强的时期。你会更加认真地对待自己和人生目标，可能会经历一些限制或挑战，但这些都是为了建立更坚实的基础。这是培养耐心、承担责任、塑造成熟人格的重要时期。",
		
		"saturn_house_7": "关系考验和成熟的时期。亲密关系面临现实的考验，需要认真对待承诺和责任。这可能带来关系的巩固或是清理不健康的连接。适合建立长期稳定的伙伴关系，学习关系中的责任和界限。",
		
		"saturn_house_10": "事业建构的关键时期。你会面对职业上的重大责任和挑战，可能需要付出更多努力。这是建立职业声誉、实现长期职业目标的重要阶段。成功需要耐心、纪律和持之以恒的努力。",
	}
	
	if text, ok := interpretations[key]; ok {
		return text
	}
	
	// House-specific default interpretations (维度相关)
	houseDefaults := map[string]string{
		"1":  "这段时间你的**个人活力**和**职业表现**都会受到影响。关注**身体状态**如何影响**工作效率**，保持良好的**精力管理**对**事业发展**很重要。",
		"2":  "这段时间**财务状况**成为关注点。留意**收入变化**、**资产管理**和**消费习惯**，这些都会影响你的**物质安全感**。",
		"3":  "这段时间**沟通交流**会带来**职业机会**和**人际连接**。注意**工作中的合作**和**社交网络**的拓展，两者相辅相成。",
		"4":  "这段时间**家庭关系**和**情感纽带**成为焦点。关注与**家人的相处**、**亲密关系的质量**，**情感安全**是这个阶段的主题。",
		"5":  "这段时间**浪漫爱情**和**创造表达**交织在一起。通过**艺术创作**和**娱乐活动**来**展现真实自我**，在**亲密关系**中也更加真诚。",
		"6":  "这段时间**健康状况**和**工作表现**密切相关。良好的**身体管理**支持**职业效率**，**规律作息**带来**工作稳定性**。",
		"7":  "这段时间**人际关系**和**合作伙伴**成为重点。关注**一对一的深度连接**、**婚姻感情**和**商业合作**等**关系议题**。",
		"8":  "这段时间**内在转化**和**财务共享**同时发生。**心理成长**伴随着**共同资产**的处理，**灵性蜕变**与**物质整合**并行。",
		"9":  "这段时间**精神追求**和**智慧探索**成为主题。关注**哲学思考**、**信仰体系**、**高等教育**等**灵性成长**的议题。",
		"10": "这段时间**职业发展**和**事业成就**是核心。关注**工作表现**、**职位晋升**、**专业声誉**等**事业相关**的重要事项。",
		"11": "这段时间**社交网络**和**共同理想**成为焦点。通过**友谊圈子**获得**精神启发**，在**集体活动**中**提升意识**。",
		"12": "这段时间**灵性修行**和**内在疗愈**是主题。关注**冥想静心**、**潜意识探索**、**心灵净化**等**精神层面**的成长。",
	}
	
	if defaultText, ok := houseDefaults[house]; ok {
		return defaultText
	}
	
	// Final fallback
	return "这段时间你会在这个生活领域经历重要的发展和变化。注意观察相关的事件和感受，它们会为你的人生带来重要的启示。"
}

// ========== English Transit House Interpretations ==========

func getEnglishTransitHouseInterpretation(key string, house string) string {
	interpretations := map[string]string{
		// 1st House: Health + Career (dual)
		"sun_house_1": "Personal charisma and **vital energy** are fully enhanced. You'll feel **physically energized** with excellent mental state, **stamina and endurance** at their peak. This powerful life force makes you **excel in professional settings**, with **leadership and personal magnetism** attracting attention. Good time to take initiative at **work**, showcasing your **professional competence and healthy vitality**. Maintain good **rest habits** to let abundant energy support **career development**.",
		
		// 2nd House: Finance (single)
		"sun_house_2": "**Financial matters and material security** become the focus. You'll pay more attention to **income sources** and **asset status**, with possible opportunities to **increase earnings** or reassess **financial planning**. Clearer understanding of your value enhances **earning capacity**. Good time for **making financial plans, finding new income sources**, or **investing in long-term growth**. Review personal **money values** to determine what **material foundation** truly matters to you.",
		
		// 3rd House: Career + Relationship (dual)
		"sun_house_3": "Communication brings **career opportunities** and **interpersonal expansion**. You'll find **workplace collaboration** becoming frequent with more **business meetings and presentations**. Meanwhile, **social networks** expand as **making new friends and maintaining old relationships** goes smoothly. Good time for **building professional connections** and **advancing projects through communication skills**. Short **business trips** may bring **commercial cooperation**, while **social activities** also support **career growth**. Balance between **work communication** and **personal socializing**.",
		
		// 4th House: Relationship (single)
		"sun_house_4": "**Family relationships and emotional bonds** become life's center. You'll focus more on **time with family**, desiring to establish or improve living environment, with **quality of intimate relationships** becoming more important. Good period for **repairing family ties and deepening emotional connections**. You may review how **childhood family patterns** affect current **intimate relationships**. Suitable for spending more time **with family, deepening emotional bonds**, or **improving home environment** to make it warmer.",
		
		// 10th House: Career (single)
		"sun_house_10": "**Career development** and **professional achievement** become the center of attention. Your **professional path** and **social status** are emphasized, with possible **promotion opportunities** or more **work responsibilities**. Your **professional abilities** and **performance results** gain public recognition, ideal for **setting career goals and advancing key projects**. Relationships with **superiors** or **industry authorities** become important. Critical period for building **professional reputation** and realizing **career aspirations**. Showcase **professional excellence** and **pursue excellence at work**, while maintaining humility and team respect.",
		
		"moon_house_1": "You may feel stronger emotional fluctuations these days. You might become more expressive, with external behavior influenced by emotions. You'll find yourself more sensitive to the outside world and more easily affected by others' emotions. This period is suitable for expressing genuine feelings, but pay attention to emotional management to avoid overreacting.",
		
		"venus_house_11": "Social charisma blooms as friendships and community activities bring joy. You'll make new friends and participate in interesting group activities. Your social network expands, and you might meet romantic interests through friends. This is a great time to engage in social activities, pursue collective goals, and enjoy the beauty of friendship.",
		
		"mars_house_11": "Bold and brave pursuit of ideals and social goals. You'll participate in group activities, fighting for common objectives. You show drive and leadership in social settings. This period is suitable for promoting social change and achieving collective goals, but remember to cooperate with the team.",
		
		"jupiter_house_11": "Benefactors appear - seize the opportunity. Friendships and social activities bring luck and opportunities. You'll meet influential friends, expanding your social network. Participating in group activities brings growth and rewards. This is a good period for achieving long-term goals and gaining social support.",
	}
	
	if text, ok := interpretations[key]; ok {
		return text
	}
	
	// House-specific default interpretations (dimension-focused)
	houseDefaults := map[string]string{
		"1":  "Your **personal vitality** and **professional performance** are affected this period. Notice how **physical state** impacts **work efficiency**; good **energy management** is important for **career development**.",
		"2":  "**Financial matters** become the focus. Pay attention to **income changes**, **asset management**, and **spending habits** affecting your **material security**.",
		"3":  "**Communication** brings **career opportunities** and **interpersonal connections**. Notice **workplace collaboration** and **social network** expansion working together.",
		"4":  "**Family relationships** and **emotional bonds** are highlighted. Focus on **time with family** and **quality of intimate relationships**; **emotional security** is the theme.",
		"5":  "**Romantic love** and **creative expression** intertwine. Express **authentic self** through **artistic creation** and **entertainment**, being more genuine in **intimate relationships**.",
		"6":  "**Health status** and **work performance** are closely related. Good **physical management** supports **professional efficiency**; **regular routines** bring **work stability**.",
		"7":  "**Interpersonal relationships** and **partnerships** are key. Focus on **one-on-one deep connections**, **marriage**, and **business cooperation** and other **relationship matters**.",
		"8":  "**Inner transformation** and **financial sharing** occur simultaneously. **Psychological growth** accompanies **shared assets** handling; **spiritual evolution** parallels **material integration**.",
		"9":  "**Spiritual pursuits** and **wisdom exploration** are the theme. Focus on **philosophical thinking**, **belief systems**, **higher education** and other **spiritual growth** matters.",
		"10": "**Career development** and **professional achievement** are core. Focus on **work performance**, **position advancement**, **professional reputation** and other **career-related** priorities.",
		"11": "**Social networks** and **shared ideals** are highlighted. Gain **spiritual inspiration** through **friend circles**; **raise consciousness** in **collective activities**.",
		"12": "**Spiritual practice** and **inner healing** are themes. Focus on **meditation**, **subconscious exploration**, **soul purification** and other **spiritual** growth.",
	}
	
	if defaultText, ok := houseDefaults[house]; ok {
		return defaultText
	}
	
	return "During this period, you'll experience important developments and changes in this life area. Pay attention to related events and feelings, as they will bring important insights to your life."
}

// ========== Russian Transit House Interpretations ==========

func getRussianTransitHouseInterpretation(key string, house string) string {
	interpretations := map[string]string{
		// 1st House: Health + Career
		"sun_house_1": "Личное обаяние и **жизненная энергия** полностью усилены. Вы почувствуете себя **физически энергичным** с отличным ментальным состоянием, **выносливость и сила** на пике. Эта мощная жизненная сила позволяет вам **преуспевать в профессиональной среде**, **лидерские качества и личный магнетизм** привлекают внимание. Хорошее время для проявления инициативы на **работе**, демонстрируя вашу **профессиональную компетентность и здоровую жизнеспособность**.",
		
		// 2nd House: Finance
		"sun_house_2": "**Финансовые вопросы и материальная безопасность** становятся фокусом. Вы уделите больше внимания **источникам дохода** и **статусу активов**, с возможными возможностями **увеличить заработок** или пересмотреть **финансовое планирование**. Более четкое понимание вашей ценности повышает **способность зарабатывать**.",
		
		// 3rd House: Career + Relationship
		"sun_house_3": "Общение приносит **карьерные возможности** и **расширение межличностных связей**. Вы обнаружите, что **сотрудничество на рабочем месте** становится частым с большим количеством **деловых встреч и презентаций**. Между тем **социальные сети** расширяются, так как **заводить новых друзей и поддерживать старые отношения** идет гладко. Хорошее время для **построения профессиональных связей** и **продвижения проектов через коммуникативные навыки**.",
		
		// 10th House: Career
		"sun_house_10": "**Развитие карьеры** и **профессиональные достижения** становятся центром внимания. Ваш **профессиональный путь** и **социальный статус** подчеркнуты, возможны **возможности повышения** или больше **рабочих обязанностей**. Ваши **профессиональные способности** и **результаты работы** получают общественное признание. Критический период для построения **профессиональной репутации** и реализации **карьерных амбиций**.",
		
		"venus_house_11": "Расцвет социальной харизмы - дружба и общественная деятельность приносят радость. Вы заведете новых друзей, будете участвовать в интересных групповых мероприятиях. Ваша социальная сеть расширяется. Отличное время для участия в социальных мероприятиях и достижения коллективных целей.",
		
		"jupiter_house_11": "Появляются покровители - ловите возможность. Дружба и социальная активность приносят удачу и возможности. Вы познакомитесь с влиятельными друзьями, расширите свою социальную сеть. Участие в групповых мероприятиях принесет рост и награды.",
	}
	
	if text, ok := interpretations[key]; ok {
		return text
	}
	
	// House-specific defaults (dimension-focused)
	houseDefaults := map[string]string{
		"1":  "Ваша **личная жизненная сила** и **профессиональная деятельность** затронуты. Замечайте, как **физическое состояние** влияет на **рабочую эффективность**; хорошее **управление энергией** важно для **развития карьеры**.",
		"2":  "**Финансовые вопросы** становятся фокусом. Обратите внимание на **изменения дохода**, **управление активами** и **привычки расходов**, влияющие на вашу **материальную безопасность**.",
		"3":  "**Общение** приносит **карьерные возможности** и **межличностные связи**. Замечайте **сотрудничество на работе** и **расширение социальной сети**, работающие вместе.",
		"4":  "**Семейные отношения** и **эмоциональные связи** выделены. Сосредоточьтесь на **времени с семьей** и **качестве интимных отношений**; **эмоциональная безопасность** является темой.",
		"10": "**Развитие карьеры** и **профессиональные достижения** являются основными. Сосредоточьтесь на **рабочих результатах**, **продвижении по службе**, **профессиональной репутации**.",
	}
	
	if defaultText, ok := houseDefaults[house]; ok {
		return defaultText
	}
	
	return "В этот период вы переживете важные события и изменения в этой сфере жизни. Обратите внимание на связанные события и чувства, они принесут важные откровения в вашу жизнь."
}

// ========== Progression Interpretations ==========

func getChineseProgressionInterpretation(key string, isPositive bool) string {
	interpretations := map[string]string{
		"moon_trine_venus": "这是情感和谐、关系美好的发展阶段。你会感到内心平和，容易与他人建立温暖的连接。艺术品味提升，对美的事物特别敏感。这段时间适合深化亲密关系，享受情感的满足。你的魅力和吸引力自然流露，人际关系顺畅愉快。",
		
		"moon_square_saturn": "这是情感成熟但也可能感到压抑的时期。你会面对现实的限制，需要在情感需求和责任之间寻求平衡。可能会经历孤独或情感的挑战，但这些经历会让你变得更加成熟和稳定。适合建立情感边界，学习自我关怀。",
		
		"sun_sextile_sun": "这是自我认同和生命力增强的时期。你会更清楚地了解自己是谁，想要什么。自信心自然提升，能够更好地表达真实的自己。这段时间适合追求个人目标，发展个性，展现领导能力。你与自己的核心本质更加协调一致。",
		
		"mercury_sextile_venus": "这是优雅沟通和思维和谐的时期。你的表达变得更加温和、富有魅力，容易与他人建立愉快的交流。艺术和美学的理解力提升，适合从事创意写作、设计等工作。人际沟通顺畅，容易获得他人的好感和支持。",
		
		"venus_conjunction_sun": "这是魅力绽放、爱与自我融合的重要时期。你的个人魅力达到高峰，容易吸引他人注意和喜爱。艺术才能和审美能力增强，适合追求浪漫、从事艺术创作。自我价值感提升，你会更加欣赏和爱护自己。",
	}
	
	if text, ok := interpretations[key]; ok {
		return text
	}
	
	if isPositive {
		return "这是积极发展和个人成长的时期。你会在这个领域经历和谐的进展，内在能力得到提升。适合把握机会，顺应这股发展的能量。"
	}
	return "这是面对挑战、促进成长的时期。你会在这个领域经历一些困难或调整，但这些经历会带来深层的成熟和理解。保持耐心，从经验中学习。"
}

func getEnglishProgressionInterpretation(key string, isPositive bool) string {
	interpretations := map[string]string{
		"moon_trine_venus": "This is a phase of emotional harmony and beautiful relationships. You'll feel inner peace and easily establish warm connections with others. Artistic taste improves, and you're particularly sensitive to beauty. This period is ideal for deepening intimate relationships and enjoying emotional fulfillment.",
		
		"sun_sextile_sun": "This is a period of enhanced self-recognition and vitality. You'll have a clearer understanding of who you are and what you want. Confidence naturally increases, allowing better expression of your true self. This time is suitable for pursuing personal goals and developing your individuality.",
		
		"venus_conjunction_sun": "This is an important period where charm blooms and love merges with self. Your personal charisma peaks, easily attracting others' attention and affection. Artistic talent and aesthetic ability are enhanced. Self-worth increases as you appreciate and love yourself more.",
	}
	
	if text, ok := interpretations[key]; ok {
		return text
	}
	
	if isPositive {
		return "This is a period of positive development and personal growth. You'll experience harmonious progress in this area, with inner abilities being enhanced."
	}
	return "This is a period of facing challenges and promoting growth. You'll experience some difficulties or adjustments in this area, but these experiences will bring deep maturity and understanding."
}

func getRussianProgressionInterpretation(key string, isPositive bool) string {
	interpretations := map[string]string{
		"sun_sextile_sun": "Это период усиленного самопознания и жизненной силы. У вас будет более четкое понимание того, кто вы есть и чего хотите. Уверенность естественно возрастает, позволяя лучше выражать свое истинное я.",
	}
	
	if text, ok := interpretations[key]; ok {
		return text
	}
	
	if isPositive {
		return "Это период позитивного развития и личностного роста. Вы испытаете гармоничный прогресс в этой области с усилением внутренних способностей."
	}
	return "Это период столкновения с вызовами и содействия росту. Вы столкнетесь с некоторыми трудностями в этой области, но эти переживания принесут глубокую зрелость и понимание."
}

// ========== Aspect Interpretations ==========

func getChineseAspectInterpretation(key string, isPositive bool) string {
	interpretations := map[string]string{
		// ========== 自合相 (Conjunction with Self) ==========
		// 太阳自合相 - 年度重要周期
		"sun_conjunction_sun": "这是你**个人身份和生命力**的重要更新周期。行运太阳回到你本命太阳的位置，标志着**新一年的开始**和**个人目标的重新设定**。你会感到**能量充沛**，**自信心增强**，对**未来方向**有了更清晰的认识。这是**设定年度目标**、**展现领导力**、**追求个人理想**的绝佳时机。你的**核心本质**得到强化，**自我表达**更加自然。适合**启动新项目**、**展现专业能力**，在**职业和生活中**追求**个人成就**。",
		
		// 月亮自合相 - 月度情感周期
		"moon_conjunction_moon": "这是**情感和内在需求**的回归时刻。行运月亮回到你本命月亮的位置，带来**情感上的熟悉感**和**安全感**。你会重新连接到**内心深处的感受**，对**情感需求**有了更清晰的认识。这是**情感疗愈**、**家庭团聚**、**内心平静**的好时机。你的**直觉力**增强，能够**理解自己的情绪**并做出**符合内心**的决定。适合**照顾自己**、**与家人相处**、**处理情感议题**，让**内在的滋养**支持你的**日常生活**。",
		
		// 水星自合相 - 思维更新
		"mercury_conjunction_mercury": "这是**思维和沟通**的重要更新期。行运水星回到你本命水星的位置，带来**思维上的清晰**和**表达能力的提升**。你会感到**学习能力增强**，对**信息处理**更加高效，**沟通交流**更加顺畅。这是**学习新知识**、**写作创作**、**商务洽谈**的好时机。你的**逻辑思维**和**分析能力**得到强化，适合**制定计划**、**解决问题**、**与他人交流想法**。注意**信息过载**，保持**专注**和**理性思考**。",
		
		// 金星自合相 - 爱与美的回归
		"venus_conjunction_venus": "这是**爱情和美感**的重要周期。行运金星回到你本命金星的位置，带来**魅力的提升**和**关系的和谐**。你会感到**个人魅力**增强，**审美能力**提升，对**美好事物**更加敏感。这是**追求浪漫**、**深化感情**、**享受艺术**的好时机。你的**人际关系**更加顺畅，**社交活动**带来愉悦。适合**表达爱意**、**改善关系**、**从事创意工作**，让**爱与美**充满你的生活。",
		
		// 火星自合相 - 行动力回归
		"mars_conjunction_mars": "这是**行动力和勇气**的重要更新期。行运火星回到你本命火星的位置，带来**能量的爆发**和**主动性的增强**。你会感到**精力充沛**，**行动力**提升，**竞争意识**增强。这是**启动新项目**、**追求目标**、**展现领导力**的好时机。你的**决断力**和**执行力**得到强化，适合**克服障碍**、**争取机会**、**展现个人能力**。注意**控制脾气**，避免**冲动行事**，将**能量**用于**建设性**的目标。",
		
		// 木星自合相 - 12年大周期
		"jupiter_conjunction_jupiter": "这是**成长和机遇**的12年大周期。行运木星回到你本命木星的位置，标志着**人生新阶段的开始**和**重要机遇的到来**。你会感到**视野开阔**，**乐观精神**增强，对**未来**充满希望。这是**扩展人生**、**追求理想**、**获得成长**的绝佳时机。你的**幸运指数**提升，**贵人运**增强，**事业发展**和**个人成长**都有重要突破。适合**设定远大目标**、**学习深造**、**拓展视野**，让**木星的祝福**引领你走向**成功**。",
		
		// 土星自合相 - 29年大周期
		"saturn_conjunction_saturn": "这是**责任和成熟**的29年大周期。行运土星回到你本命土星的位置，标志着**人生重要转折点**和**重大责任的到来**。你会感到**压力增加**，但也**更加成熟**，对**人生目标**有了更深刻的认识。这是**承担责任**、**建立结构**、**实现长期目标**的关键时期。你的**自律性**和**耐心**得到考验，**职业发展**和**人生规划**都需要**认真对待**。适合**制定长期计划**、**建立稳定基础**、**面对现实挑战**，让**土星的教导**帮助你**成长**。",
		
		// ========== 月亮相位 (Moon Aspects) ==========
		// 月亮与太阳
		"moon_conjunction_sun": "这是**情感与意志**的融合时刻。**内在感受**与**外在表达**达到和谐，你会感到**情感上的满足**和**个人身份的清晰**。这是**情感表达**、**自我展现**的好时机。",
		"moon_trine_sun": "**情感与自我**和谐共振，带来**内心的平静**和**自信的提升**。你会感到**情感需求**得到满足，**个人表达**更加自然。适合**展现真实自我**、**深化情感连接**。",
		"moon_square_sun": "**情感需求**与**个人目标**之间可能存在**冲突**。你需要在**内心感受**和**外在追求**之间找到平衡。这是**自我整合**的机会，学习**协调不同需求**。",
		"moon_sextile_sun": "**情感与自我**相互支持，带来**温和的成长**。你会感到**情感稳定**，**个人发展**顺利。适合**表达情感**、**追求个人目标**，两者**相辅相成**。",
		"moon_opposition_sun": "**情感与意志**形成**对立**，需要**平衡内在需求**和**外在表现**。你可能会感到**情感波动**，但这也是**自我认识**的机会。学习**整合**不同的**自我面向**。",
		
		// 月亮与金星
		"moon_conjunction_venus": "这是**情感与爱情**的美好融合。**内心感受**与**爱的表达**和谐统一，你会感到**情感上的满足**和**关系的和谐**。这是**深化感情**、**享受美好**的好时机。",
		"moon_trine_venus": "**情感与美感**和谐共振，带来**内心的愉悦**和**关系的顺畅**。你会感到**情感需求**得到满足，**审美能力**提升。适合**表达爱意**、**享受艺术**、**改善关系**。",
		"moon_square_venus": "**情感需求**与**爱的表达**之间可能存在**不协调**。你需要在**内心感受**和**关系期待**之间找到平衡。这是**情感成长**的机会，学习**真实表达**自己的需求。",
		"moon_sextile_venus": "**情感与美感**相互支持，带来**温和的愉悦**。你会感到**情感稳定**，**人际关系**顺畅。适合**表达情感**、**享受美好**，让**爱与美**充满生活。",
		
		// 月亮与火星
		"moon_conjunction_mars": "这是**情感与行动**的强烈融合。**内心感受**驱动**外在行动**，你会感到**情绪激动**和**行动力增强**。这是**表达情感**、**采取行动**的好时机，但要注意**控制情绪**。",
		"moon_trine_mars": "**情感与行动**和谐共振，带来**内在动力**和**外在表现**的统一。你会感到**情感能量**转化为**行动力**，**目标追求**更加顺畅。适合**主动出击**、**追求目标**。",
		"moon_square_mars": "**情感与行动**之间可能存在**冲突**。你可能会感到**情绪激动**或**行动冲动**。这是**情绪管理**的机会，学习**理性行动**，避免**冲动决定**。",
		"moon_sextile_mars": "**情感与行动**相互支持，带来**温和的动力**。你会感到**情感稳定**，**行动力**适中。适合**表达情感**、**采取行动**，两者**平衡发展**。",
		
		// 月亮与木星
		"moon_conjunction_jupiter": "这是**情感与成长**的美好融合。**内心感受**与**乐观精神**和谐统一，你会感到**情感上的满足**和**视野的开阔**。这是**情感成长**、**追求理想**的好时机。",
		"moon_trine_jupiter": "**情感与乐观**和谐共振，带来**内心的喜悦**和**生活的丰富**。你会感到**情感需求**得到满足，**人生视野**更加开阔。适合**表达情感**、**追求成长**、**享受生活**。",
		"moon_square_jupiter": "**情感需求**与**理想追求**之间可能需要**平衡**。你需要在**内心感受**和**外在目标**之间找到协调。这是**情感成长**的机会，学习**真实表达**自己的需求。",
		"moon_sextile_jupiter": "**情感与成长**相互支持，带来**温和的扩展**。你会感到**情感稳定**，**人生视野**开阔。适合**表达情感**、**追求理想**，让**成长**充满生活。",
		
		// 月亮与土星
		"moon_conjunction_saturn": "这是**情感与责任**的严肃融合。**内心感受**与**现实责任**交织，你会感到**情感上的压力**但也**更加成熟**。这是**情感成长**、**承担责任**的好时机。",
		"moon_trine_saturn": "**情感与责任**和谐共振，带来**内心的稳定**和**现实的支撑**。你会感到**情感需求**得到满足，**责任感**增强。适合**建立稳定关系**、**承担责任**、**情感成熟**。",
		"moon_square_saturn": "近期是生机勃勃、精力旺盛的时期，你也许会通过卖力工作来取得不错的成绩，也非常适合启动新的项目。你的积极上进会令你获得上司的赏识和下属的尊重，但也需要警惕过度劳累。这段时间你可能会承担额外的责任，感到情感上的压力。学习在责任和自我关怀之间找到平衡，这是走向成熟的必经之路。",
		"moon_sextile_saturn": "**情感与责任**相互支持，带来**温和的稳定**。你会感到**情感稳定**，**责任感**适中。适合**表达情感**、**承担责任**，两者**平衡发展**。",
		
		// ========== 太阳相位 (Sun Aspects) ==========
		// 太阳与木星
		"sun_conjunction_jupiter": "这是**自我与成长**的完美融合。**个人身份**与**乐观精神**和谐统一，你会感到**自信提升**和**视野开阔**。这是**追求理想**、**展现领导力**的好时机。",
		"sun_trine_jupiter": "**自我与乐观**和谐共振，带来**自信的提升**和**机遇的增加**。你会感到**个人魅力**增强，**人生视野**更加开阔。适合**追求目标**、**展现能力**、**获得成长**。",
		"sun_square_jupiter": "**个人目标**与**理想追求**之间可能需要**平衡**。你需要在**现实行动**和**远大理想**之间找到协调。这是**目标设定**的机会，学习**务实追求**理想。",
		"sun_sextile_jupiter": "**自我与成长**相互支持，带来**温和的扩展**。你会感到**自信稳定**，**人生视野**开阔。适合**追求目标**、**展现能力**，让**成长**充满生活。",
		
		// 太阳与土星
		"sun_conjunction_saturn": "这是**自我与责任**的严肃融合。**个人身份**与**现实责任**交织，你会感到**压力增加**但也**更加成熟**。这是**承担责任**、**建立结构**的好时机。",
		"sun_trine_saturn": "**自我与责任**和谐共振，带来**稳定的成长**和**现实的支撑**。你会感到**个人能力**增强，**责任感**提升。适合**建立目标**、**承担责任**、**追求成就**。",
		"sun_square_saturn": "**个人目标**与**现实责任**之间可能存在**冲突**。你可能会感到**压力增加**或**限制增多**。这是**成熟成长**的机会，学习**面对现实**，**承担责任**。",
		"sun_sextile_saturn": "**自我与责任**相互支持，带来**温和的稳定**。你会感到**自信稳定**，**责任感**适中。适合**追求目标**、**承担责任**，两者**平衡发展**。",
		
		// 太阳与金星
		"sun_conjunction_venus": "这是**自我与爱情**的美好融合。**个人魅力**与**爱的表达**和谐统一，你会感到**自信提升**和**关系和谐**。这是**展现魅力**、**追求爱情**的好时机。",
		"sun_trine_venus": "**自我与美感**和谐共振，带来**魅力的提升**和**关系的顺畅**。你会感到**个人魅力**增强，**人际关系**更加和谐。适合**表达自我**、**改善关系**、**享受美好**。",
		"sun_square_venus": "**个人表达**与**关系需求**之间可能需要**平衡**。你需要在**自我展现**和**关系和谐**之间找到协调。这是**关系成长**的机会，学习**真实表达**自己。",
		"sun_sextile_venus": "**自我与美感**相互支持，带来**温和的魅力**。你会感到**自信稳定**，**人际关系**顺畅。适合**表达自我**、**改善关系**，让**爱与美**充满生活。",
		
		// ========== 金星相位 (Venus Aspects) ==========
		// 金星与火星
		"venus_conjunction_mars": "这是**爱情与激情**的强烈融合。**爱的表达**与**行动力**和谐统一，你会感到**魅力提升**和**关系活跃**。这是**追求爱情**、**表达激情**的好时机。",
		"venus_trine_mars": "**爱情与行动**和谐共振，带来**关系的活跃**和**魅力的提升**。你会感到**关系顺畅**，**行动力**增强。适合**追求爱情**、**表达情感**、**采取行动**。",
		"venus_square_mars": "**爱的表达**与**行动力**之间可能存在**冲突**。你可能会感到**关系紧张**或**行动冲动**。这是**关系成长**的机会，学习**理性表达**，避免**冲动决定**。",
		"venus_sextile_mars": "**爱情与行动**相互支持，带来**温和的活跃**。你会感到**关系稳定**，**行动力**适中。适合**表达情感**、**采取行动**，两者**平衡发展**。",
		
		// 金星与木星
		"venus_conjunction_jupiter": "这是**爱情与成长**的美好融合。**爱的表达**与**乐观精神**和谐统一，你会感到**关系和谐**和**视野开阔**。这是**深化感情**、**追求理想**的好时机。",
		"venus_trine_jupiter": "**爱情与乐观**和谐共振，带来**关系的顺畅**和**生活的丰富**。你会感到**关系和谐**，**人生视野**更加开阔。适合**表达爱意**、**追求成长**、**享受生活**。",
		"venus_square_jupiter": "**关系需求**与**理想追求**之间可能需要**平衡**。你需要在**关系和谐**和**个人目标**之间找到协调。这是**关系成长**的机会，学习**真实表达**自己的需求。",
		"venus_sextile_jupiter": "**爱情与成长**相互支持，带来**温和的扩展**。你会感到**关系稳定**，**人生视野**开阔。适合**表达爱意**、**追求理想**，让**成长**充满生活。",
		
		// ========== 外行星相位 (Outer Planet Aspects) ==========
		// 木星与土星
		"jupiter_conjunction_saturn": "这是**成长与责任**的重要融合。**乐观精神**与**现实责任**交织，你会感到**机遇与挑战**并存。这是**建立长期目标**、**承担责任**的好时机。",
		"jupiter_trine_saturn": "**成长与责任**和谐共振，带来**稳定的扩展**和**现实的支撑**。你会感到**机遇增加**，**责任感**提升。适合**追求目标**、**承担责任**、**实现理想**。",
		"jupiter_square_saturn": "**理想追求**与**现实责任**之间可能存在**冲突**。你可能会感到**压力增加**或**限制增多**。这是**成熟成长**的机会，学习**面对现实**，**务实追求**理想。",
		"jupiter_sextile_saturn": "**成长与责任**相互支持，带来**温和的稳定**。你会感到**机遇稳定**，**责任感**适中。适合**追求目标**、**承担责任**，两者**平衡发展**。",
		
		// 土星与天王星
		"saturn_conjunction_uranus": "这是**传统与变革**的激烈融合。**稳定结构**与**创新突破**交织，你会感到**压力与机遇**并存。这是**面对变革**、**建立新结构**的好时机。",
		"saturn_trine_uranus": "**传统与变革**和谐共振，带来**稳定的创新**和**结构的突破**。你会感到**变革顺利**，**稳定性**增强。适合**创新突破**、**建立新结构**、**追求变革**。",
		"saturn_square_uranus": "**稳定结构**与**创新突破**之间可能存在**冲突**。你可能会感到**压力增加**或**变革困难**。这是**成长突破**的机会，学习**平衡传统与创新**。",
		
		// ========== 其他重要相位 ==========
		// 水星与金星
		"mercury_conjunction_venus": "这是**思维与美感**的优雅融合。**沟通表达**与**审美能力**和谐统一，你会感到**表达流畅**和**魅力提升**。这是**创意写作**、**艺术表达**的好时机。",
		"mercury_trine_venus": "**思维与美感**和谐共振，带来**表达的优雅**和**沟通的顺畅**。你会感到**沟通能力**增强，**审美能力**提升。适合**创意表达**、**改善沟通**、**享受艺术**。",
		
		// 水星与火星
		"mercury_conjunction_mars": "这是**思维与行动**的强烈融合。**沟通表达**与**行动力**和谐统一，你会感到**表达有力**和**行动迅速**。这是**快速决策**、**有效沟通**的好时机。",
		"mercury_trine_mars": "**思维与行动**和谐共振，带来**表达的力度**和**行动的效率**。你会感到**沟通能力**增强，**行动力**提升。适合**快速决策**、**有效沟通**、**采取行动**。",
		"mercury_square_mars": "**思维与行动**之间可能存在**冲突**。你可能会感到**表达冲动**或**行动急躁**。这是**沟通成长**的机会，学习**理性表达**，避免**冲动决定**。",
	}
	
	if text, ok := interpretations[key]; ok {
		return text
	}
	
	// Fallback based on aspect type and positivity
	if isPositive {
		return "这段时间两颗行星的能量和谐共振，为你带来积极的发展机会。适合把握这股顺势的能量，在相关领域采取行动，会获得良好的结果。"
	}
	return "这段时间两颗行星形成挑战性的角度，可能会带来一些紧张或需要调整的情况。这是学习和成长的机会，面对挑战会让你变得更加成熟和强大。"
}

func getEnglishAspectInterpretation(key string, isPositive bool) string {
	interpretations := map[string]string{
		// ========== Self Conjunctions ==========
		"sun_conjunction_sun": "This is an important renewal cycle for your **personal identity and vitality**. The transiting Sun returns to your natal Sun's position, marking the **beginning of a new year** and **reset of personal goals**. You'll feel **energized**, **more confident**, with clearer understanding of **future direction**. This is an ideal time to **set annual goals**, **show leadership**, and **pursue personal ideals**. Your **core essence** is strengthened, and **self-expression** becomes more natural. Suitable for **starting new projects**, **demonstrating professional abilities**, and pursuing **personal achievements** in **career and life**.",
		
		"moon_conjunction_moon": "This is a moment of return for **emotions and inner needs**. The transiting Moon returns to your natal Moon's position, bringing **emotional familiarity** and **sense of security**. You'll reconnect with **deep inner feelings** and gain clearer understanding of **emotional needs**. This is a good time for **emotional healing**, **family reunions**, and **inner peace**. Your **intuition** strengthens, enabling you to **understand your emotions** and make decisions **aligned with your heart**. Suitable for **self-care**, **spending time with family**, and **addressing emotional issues**.",
		
		"mercury_conjunction_mercury": "This is an important renewal period for **thinking and communication**. The transiting Mercury returns to your natal Mercury's position, bringing **mental clarity** and **enhanced expression**. You'll feel **improved learning ability**, more efficient **information processing**, and smoother **communication**. This is a good time for **learning new knowledge**, **creative writing**, and **business negotiations**. Your **logical thinking** and **analytical abilities** are strengthened. Be mindful of **information overload**; maintain **focus** and **rational thinking**.",
		
		"venus_conjunction_venus": "This is an important cycle for **love and beauty**. The transiting Venus returns to your natal Venus's position, bringing **enhanced charm** and **relationship harmony**. You'll feel **increased personal magnetism**, **improved aesthetic sense**, and greater sensitivity to **beautiful things**. This is a good time to **pursue romance**, **deepen relationships**, and **enjoy art**. Your **interpersonal relationships** become smoother, and **social activities** bring joy. Suitable for **expressing love**, **improving relationships**, and **creative work**.",
		
		"mars_conjunction_mars": "This is an important renewal period for **action and courage**. The transiting Mars returns to your natal Mars's position, bringing **energy bursts** and **increased initiative**. You'll feel **energized**, **more active**, and **more competitive**. This is a good time to **start new projects**, **pursue goals**, and **show leadership**. Your **decisiveness** and **execution** are strengthened. Be mindful of **controlling temper**; avoid **impulsive actions** and channel **energy** toward **constructive** goals.",
		
		"jupiter_conjunction_jupiter": "This is a major 12-year cycle of **growth and opportunity**. The transiting Jupiter returns to your natal Jupiter's position, marking the **beginning of a new life phase** and **arrival of important opportunities**. You'll feel **expanded vision**, **increased optimism**, and **hope for the future**. This is an ideal time to **expand your life**, **pursue ideals**, and **achieve growth**. Your **luck index** rises, **benefactor connections** increase, and both **career development** and **personal growth** see important breakthroughs. Suitable for **setting ambitious goals**, **advanced learning**, and **expanding horizons**.",
		
		"saturn_conjunction_saturn": "This is a major 29-year cycle of **responsibility and maturity**. The transiting Saturn returns to your natal Saturn's position, marking a **major life turning point** and **arrival of significant responsibilities**. You'll feel **increased pressure** but also **greater maturity**, with deeper understanding of **life goals**. This is a critical period for **taking responsibility**, **building structure**, and **achieving long-term goals**. Your **self-discipline** and **patience** are tested. Suitable for **making long-term plans**, **building stable foundations**, and **facing reality**.",
		
		// ========== Moon Aspects ==========
		"moon_conjunction_sun": "This is a moment of fusion between **emotion and will**. **Inner feelings** and **outer expression** reach harmony. You'll feel **emotional fulfillment** and **clarity of personal identity**. This is a good time for **emotional expression** and **self-expression**.",
		
		"moon_trine_sun": "**Emotion and self** resonate harmoniously, bringing **inner peace** and **increased confidence**. You'll feel **emotional needs** are met, and **personal expression** becomes more natural. Suitable for **showing your true self** and **deepening emotional connections**.",
		
		"moon_square_sun": "There may be **conflict** between **emotional needs** and **personal goals**. You need to find balance between **inner feelings** and **outer pursuits**. This is an opportunity for **self-integration**; learn to **coordinate different needs**.",
		
		"moon_sextile_sun": "**Emotion and self** support each other, bringing **gentle growth**. You'll feel **emotionally stable** and **personal development** proceeds smoothly. Suitable for **expressing emotions** and **pursuing personal goals**.",
		
		"moon_conjunction_venus": "This is a beautiful fusion of **emotion and love**. **Inner feelings** and **love expression** harmonize. You'll feel **emotional fulfillment** and **relationship harmony**. This is a good time to **deepen relationships** and **enjoy beauty**.",
		
		"moon_trine_venus": "**Emotion and beauty** resonate harmoniously, bringing **inner joy** and **smooth relationships**. You'll feel **emotional needs** are met and **aesthetic sense** improves. Suitable for **expressing love**, **enjoying art**, and **improving relationships**.",
		
		"moon_square_venus": "There may be **disharmony** between **emotional needs** and **love expression**. You need to find balance between **inner feelings** and **relationship expectations**. This is an opportunity for **emotional growth**; learn to **express your needs authentically**.",
		
		"moon_conjunction_mars": "This is a strong fusion of **emotion and action**. **Inner feelings** drive **outer actions**. You'll feel **emotionally excited** and **more active**. This is a good time for **expressing emotions** and **taking action**, but be mindful of **controlling emotions**.",
		
		"moon_trine_mars": "**Emotion and action** resonate harmoniously, bringing unity of **inner drive** and **outer expression**. You'll feel **emotional energy** transforms into **action power**. Suitable for **taking initiative** and **pursuing goals**.",
		
		"moon_square_mars": "There may be **conflict** between **emotion and action**. You may feel **emotionally agitated** or **impulsive**. This is an opportunity for **emotional management**; learn to **act rationally** and avoid **impulsive decisions**.",
		
		"moon_conjunction_jupiter": "This is a beautiful fusion of **emotion and growth**. **Inner feelings** and **optimism** harmonize. You'll feel **emotional fulfillment** and **expanded vision**. This is a good time for **emotional growth** and **pursuing ideals**.",
		
		"moon_trine_jupiter": "**Emotion and optimism** resonate harmoniously, bringing **inner joy** and **life richness**. You'll feel **emotional needs** are met and **life vision** expands. Suitable for **expressing emotions**, **pursuing growth**, and **enjoying life**.",
		
		"moon_conjunction_saturn": "This is a serious fusion of **emotion and responsibility**. **Inner feelings** and **real-world responsibilities** intertwine. You'll feel **emotional pressure** but also **greater maturity**. This is a good time for **emotional growth** and **taking responsibility**.",
		
		"moon_trine_saturn": "**Emotion and responsibility** resonate harmoniously, bringing **inner stability** and **real-world support**. You'll feel **emotional needs** are met and **sense of responsibility** increases. Suitable for **building stable relationships**, **taking responsibility**, and **emotional maturity**.",
		
		"moon_square_saturn": "This is a period of vitality and vigor. You may achieve good results through hard work and it's very suitable for starting new projects. Your proactivity will earn you appreciation from superiors and respect from subordinates, but be wary of overwork. You may take on extra responsibilities and feel emotional pressure. Learning to balance duty and self-care is a necessary path to maturity.",
		
		// ========== Sun Aspects ==========
		"sun_conjunction_jupiter": "This is a perfect fusion of **self and growth**. **Personal identity** and **optimism** harmonize. You'll feel **increased confidence** and **expanded vision**. This is a good time to **pursue ideals** and **show leadership**.",
		
		"sun_trine_jupiter": "**Self and optimism** resonate harmoniously, bringing **increased confidence** and **more opportunities**. You'll feel **enhanced personal charm** and **expanded life vision**. Suitable for **pursuing goals**, **demonstrating abilities**, and **achieving growth**.",
		
		"sun_conjunction_saturn": "This is a serious fusion of **self and responsibility**. **Personal identity** and **real-world responsibilities** intertwine. You'll feel **increased pressure** but also **greater maturity**. This is a good time to **take responsibility** and **build structure**.",
		
		"sun_trine_saturn": "**Self and responsibility** resonate harmoniously, bringing **stable growth** and **real-world support**. You'll feel **enhanced personal abilities** and **increased responsibility**. Suitable for **setting goals**, **taking responsibility**, and **pursuing achievements**.",
		
		"sun_square_saturn": "There may be **conflict** between **personal goals** and **real-world responsibilities**. You may feel **increased pressure** or **more limitations**. This is an opportunity for **mature growth**; learn to **face reality** and **take responsibility**.",
		
		"sun_conjunction_venus": "This is a beautiful fusion of **self and love**. **Personal charm** and **love expression** harmonize. You'll feel **increased confidence** and **relationship harmony**. This is a good time to **show charm** and **pursue love**.",
		
		"sun_trine_venus": "**Self and beauty** resonate harmoniously, bringing **enhanced charm** and **smooth relationships**. You'll feel **enhanced personal magnetism** and **more harmonious relationships**. Suitable for **self-expression**, **improving relationships**, and **enjoying beauty**.",
		
		// ========== Venus Aspects ==========
		"venus_conjunction_mars": "This is a strong fusion of **love and passion**. **Love expression** and **action power** harmonize. You'll feel **enhanced charm** and **active relationships**. This is a good time to **pursue love** and **express passion**.",
		
		"venus_trine_mars": "**Love and action** resonate harmoniously, bringing **active relationships** and **enhanced charm**. You'll feel **smooth relationships** and **increased action power**. Suitable for **pursuing love**, **expressing emotions**, and **taking action**.",
		
		"venus_square_mars": "There may be **conflict** between **love expression** and **action power**. You may feel **relationship tension** or **impulsive actions**. This is an opportunity for **relationship growth**; learn to **express rationally** and avoid **impulsive decisions**.",
		
		"venus_conjunction_jupiter": "This is a beautiful fusion of **love and growth**. **Love expression** and **optimism** harmonize. You'll feel **relationship harmony** and **expanded vision**. This is a good time to **deepen relationships** and **pursue ideals**.",
		
		"venus_trine_jupiter": "**Love and optimism** resonate harmoniously, bringing **smooth relationships** and **life richness**. You'll feel **relationship harmony** and **expanded life vision**. Suitable for **expressing love**, **pursuing growth**, and **enjoying life**.",
		
		// ========== Outer Planet Aspects ==========
		"jupiter_conjunction_saturn": "This is an important fusion of **growth and responsibility**. **Optimism** and **real-world responsibilities** intertwine. You'll feel **opportunities and challenges** coexist. This is a good time to **set long-term goals** and **take responsibility**.",
		
		"jupiter_trine_saturn": "**Growth and responsibility** resonate harmoniously, bringing **stable expansion** and **real-world support**. You'll feel **increased opportunities** and **enhanced responsibility**. Suitable for **pursuing goals**, **taking responsibility**, and **achieving ideals**.",
		
		"jupiter_square_saturn": "There may be **conflict** between **ideal pursuits** and **real-world responsibilities**. You may feel **increased pressure** or **more limitations**. This is an opportunity for **mature growth**; learn to **face reality** and **pursue ideals practically**.",
		
		"saturn_conjunction_uranus": "This is an intense fusion of **tradition and change**. **Stable structure** and **innovative breakthroughs** intertwine. You'll feel **pressure and opportunities** coexist. This is a good time to **face change** and **build new structures**.",
		
		"saturn_trine_uranus": "**Tradition and change** resonate harmoniously, bringing **stable innovation** and **structural breakthroughs**. You'll feel **smooth change** and **enhanced stability**. Suitable for **innovative breakthroughs**, **building new structures**, and **pursuing change**.",
		
		"saturn_square_uranus": "There may be **conflict** between **stable structure** and **innovative breakthroughs**. You may feel **increased pressure** or **difficult change**. This is an opportunity for **growth breakthrough**; learn to **balance tradition and innovation**.",
		
		// ========== Other Important Aspects ==========
		"mercury_conjunction_venus": "This is an elegant fusion of **thinking and beauty**. **Communication** and **aesthetic sense** harmonize. You'll feel **smooth expression** and **enhanced charm**. This is a good time for **creative writing** and **artistic expression**.",
		
		"mercury_trine_venus": "**Thinking and beauty** resonate harmoniously, bringing **elegant expression** and **smooth communication**. You'll feel **enhanced communication abilities** and **improved aesthetic sense**. Suitable for **creative expression**, **improving communication**, and **enjoying art**.",
		
		"mercury_conjunction_mars": "This is a strong fusion of **thinking and action**. **Communication** and **action power** harmonize. You'll feel **powerful expression** and **rapid action**. This is a good time for **quick decisions** and **effective communication**.",
		
		"mercury_trine_mars": "**Thinking and action** resonate harmoniously, bringing **powerful expression** and **efficient action**. You'll feel **enhanced communication abilities** and **increased action power**. Suitable for **quick decisions**, **effective communication**, and **taking action**.",
		
		"mercury_square_mars": "There may be **conflict** between **thinking and action**. You may feel **impulsive expression** or **impatient action**. This is an opportunity for **communication growth**; learn to **express rationally** and avoid **impulsive decisions**.",
	}
	
	if text, ok := interpretations[key]; ok {
		return text
	}
	
	if isPositive {
		return "During this time, the energies of these planets harmonize, bringing positive opportunities for development. It's suitable to harness this favorable energy and take action in related areas for good results."
	}
	return "During this time, these planets form a challenging angle, which may bring tension or situations requiring adjustment. This is an opportunity for learning and growth; facing challenges will make you more mature and stronger."
}

func getRussianAspectInterpretation(key string, isPositive bool) string {
	interpretations := map[string]string{
		// ========== Self Conjunctions ==========
		"sun_conjunction_sun": "Это важный цикл обновления для вашей **личной идентичности и жизненной силы**. Транзитный Солнце возвращается к положению натального Солнца, отмечая **начало нового года** и **переустановку личных целей**. Вы почувствуете себя **энергичным**, **более уверенным**, с более четким пониманием **будущего направления**. Это идеальное время для **постановки годовых целей**, **проявления лидерства** и **пursuit личных идеалов**.",
		
		"moon_conjunction_moon": "Это момент возврата для **эмоций и внутренних потребностей**. Транзитная Луна возвращается к положению натальной Луны, принося **эмоциональную знакомость** и **чувство безопасности**. Вы восстановите связь с **глубокими внутренними чувствами** и получите более четкое понимание **эмоциональных потребностей**.",
		
		"jupiter_conjunction_jupiter": "Это крупный 12-летний цикл **роста и возможностей**. Транзитный Юпитер возвращается к положению натального Юпитера, отмечая **начало новой жизненной фазы** и **приход важных возможностей**. Вы почувствуете **расширенное видение**, **возросший оптимизм** и **надежду на будущее**.",
		
		"saturn_conjunction_saturn": "Это крупный 29-летний цикл **ответственности и зрелости**. Транзитный Сатурн возвращается к положению натального Сатурна, отмечая **важный жизненный поворотный момент** и **приход значительных обязанностей**. Вы почувствуете **возросшее давление**, но также **большую зрелость**.",
		
		// ========== Moon Aspects ==========
		"moon_trine_venus": "**Эмоция и красота** резонируют гармонично, принося **внутреннюю радость** и **гладкие отношения**. Вы почувствуете, что **эмоциональные потребности** удовлетворены, а **эстетическое чувство** улучшается.",
		
		"moon_square_saturn": "Это период жизненной силы и энергии. Вы можете достичь хороших результатов через упорный труд, и это очень подходящее время для начала новых проектов. Ваша проактивность принесет вам признание начальства и уважение подчиненных, но остерегайтесь переутомления.",
		
		// ========== Sun Aspects ==========
		"sun_trine_jupiter": "**Я и оптимизм** резонируют гармонично, принося **возросшую уверенность** и **больше возможностей**. Вы почувствуете **усиленную личную привлекательность** и **расширенное жизненное видение**.",
		
		"sun_trine_saturn": "**Я и ответственность** резонируют гармонично, принося **стабильный рост** и **реальную поддержку**. Вы почувствуете **усиленные личные способности** и **возросшую ответственность**.",
	}
	
	if text, ok := interpretations[key]; ok {
		return text
	}
	
	if isPositive {
		return "В это время энергии планет гармонизируются, принося позитивные возможности для развития. Подходящее время для использования этой благоприятной энергии и действий в связанных областях."
	}
	return "В это время планеты образуют сложный угол, который может принести напряжение или ситуации, требующие корректировки. Это возможность для обучения и роста."
}

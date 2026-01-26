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
		"moon_square_saturn": "近期是生机勃勃、精力旺盛的时期，你也许会通过卖力工作来取得不错的成绩，也非常适合启动新的项目。你的积极上进会令你获得上司的赏识和下属的尊重，但也需要警惕过度劳累。这段时间你可能会承担额外的责任，感到情感上的压力。学习在责任和自我关怀之间找到平衡，这是走向成熟的必经之路。",
	}
	
	if text, ok := interpretations[key]; ok {
		return text
	}
	
	if isPositive {
		return "这段时间两颗行星的能量和谐共振，为你带来积极的发展机会。适合把握这股顺势的能量，在相关领域采取行动，会获得良好的结果。"
	}
	return "这段时间两颗行星形成挑战性的角度，可能会带来一些紧张或需要调整的情况。这是学习和成长的机会，面对挑战会让你变得更加成熟和强大。"
}

func getEnglishAspectInterpretation(key string, isPositive bool) string {
	interpretations := map[string]string{
		"moon_square_saturn": "This is a period of vitality and vigor. You may achieve good results through hard work and it's very suitable for starting new projects. Your proactivity will earn you appreciation from superiors and respect from subordinates, but be wary of overwork. You may take on extra responsibilities and feel emotional pressure. Learning to balance duty and self-care is a necessary path to maturity.",
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
	if isPositive {
		return "В это время энергии планет гармонизируются, принося позитивные возможности для развития. Подходящее время для использования этой благоприятной энергии и действий в связанных областях."
	}
	return "В это время планеты образуют сложный угол, который может принести напряжение или ситуации, требующие корректировки. Это возможность для обучения и роста."
}

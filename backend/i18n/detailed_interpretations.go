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
	
	// For retrograde events
	if eventType == "retrograde" {
		return t.getRetrogradeInterpretation(planet1)
	}
	
	// For annual lord (profection) events
	if eventType == "profectionLord" {
		return t.getProfectionLordInterpretation(planet1)
	}
	
	// For planetary hour (hourly); API may send "planetary_hour_change" or "planetaryHour"
	if eventType == "planetaryHour" || eventType == "planetary_hour_change" {
		return t.getPlanetaryHourInterpretation(planet1)
	}
	
	// For Moon void of course (hourly)
	if eventType == "voidOfCourse" {
		return t.getVoidOfCourseInterpretation()
	}
	
	// For lunar phase events
	if eventType == "lunar_phase" {
		return t.getLunarPhaseInterpretation(aspect) // aspect contains phase type
	}
	
	// For sign change events
	if eventType == "sign_change" {
		return t.getSignChangeInterpretation(planet1, aspect) // aspect contains new sign
	}
	
	// For dignity events
	if eventType == "dignity" {
		return t.getDignityInterpretation(planet1, aspect) // aspect contains dignity type
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

// getRetrogradeInterpretation returns detailed interpretation for retrograde events (planet1 = retrograde planet)
func (t *Translator) getRetrogradeInterpretation(planet models.PlanetID) string {
	switch t.lang {
	case Chinese:
		return getChineseRetrogradeInterpretation(planet)
	case Russian:
		return getRussianRetrogradeInterpretation(planet)
	default:
		return getEnglishRetrogradeInterpretation(planet)
	}
}

// getProfectionLordInterpretation returns detailed interpretation for annual lord events (planet1 = lord planet)
func (t *Translator) getProfectionLordInterpretation(planet models.PlanetID) string {
	switch t.lang {
	case Chinese:
		return getChineseProfectionLordInterpretation(planet)
	case Russian:
		return getRussianProfectionLordInterpretation(planet)
	default:
		return getEnglishProfectionLordInterpretation(planet)
	}
}

// getPlanetaryHourInterpretation returns detailed interpretation for planetary hour (planet1 = hour ruler)
func (t *Translator) getPlanetaryHourInterpretation(planet models.PlanetID) string {
	switch t.lang {
	case Chinese:
		return getChinesePlanetaryHourInterpretation(planet)
	case Russian:
		return getRussianPlanetaryHourInterpretation(planet)
	default:
		return getEnglishPlanetaryHourInterpretation(planet)
	}
}

// getVoidOfCourseInterpretation returns detailed interpretation for Moon void of course (no planet param)
func (t *Translator) getVoidOfCourseInterpretation() string {
	switch t.lang {
	case Chinese:
		return getChineseVoidOfCourseInterpretation()
	case Russian:
		return getRussianVoidOfCourseInterpretation()
	default:
		return getEnglishVoidOfCourseInterpretation()
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
		
		// Saturn through houses - key ones
		"saturn_house_1": "自律和责任感增强的时期。你会更加认真地对待自己和人生目标，可能会经历一些限制或挑战，但这些都是为了建立更坚实的基础。这是培养耐心、承担责任、塑造成熟人格的重要时期。",
		"saturn_house_2": "**财务责任**和**物质管理**成为焦点。你需要**认真规划财务**、**建立稳定的收入来源**，可能面临**财务限制**或需要**节制消费**。这是**建立长期财务安全**、**学习价值管理**的好时期。适合**制定预算**、**还清债务**、**投资于保值资产**。通过**纪律和耐心**建立**物质基础**。",
		"saturn_house_3": "**沟通责任**和**学习纪律**被强调。你需要**认真对待学习**、**改善沟通技巧**，可能感到**表达受限**或需要**更谨慎地说话**。这是**建立知识基础**、**改善与兄弟姐妹关系**的好时期。适合**系统化学习**、**完成学业**、**处理邻里关系**。通过**耐心和坚持**提升**沟通能力**。",
		"saturn_house_4": "**家庭责任**和**情感成熟**成为主题。你需要**承担家庭责任**、**处理家庭事务**，可能面临**家庭压力**或需要**建立家庭结构**。这是**改善家庭关系**、**处理房产**、**建立情感安全感**的好时期。适合**照顾家人**、**处理家庭遗产**、**建立家庭传统**。通过**责任和承诺**建立**稳定的家庭基础**。",
		"saturn_house_5": "**创造责任**和**情感纪律**被强调。你需要**认真对待创作**、**在关系中承担责任**，可能感到**创造力受限**或需要**更成熟地表达情感**。这是**建立创作习惯**、**处理与孩子的关系**、**学习情感边界**的好时期。适合**完成创作项目**、**教育孩子**、**建立健康的娱乐习惯**。",
		"saturn_house_6": "**健康责任**和**工作纪律**成为重点。你需要**认真对待健康**、**改善工作习惯**，可能面临**健康挑战**或需要**更规律的生活**。这是**建立健康习惯**、**提升职业技能**、**处理工作责任**的好时期。适合**体检**、**制定健康计划**、**改善工作流程**。通过**纪律和坚持**建立**健康的工作生活平衡**。",
		"saturn_house_7": "关系考验和成熟的时期。亲密关系面临现实的考验，需要认真对待承诺和责任。这可能带来关系的巩固或是清理不健康的连接。适合建立长期稳定的伙伴关系，学习关系中的责任和界限。",
		"saturn_house_8": "**深层责任**和**资源管理**被强调。你需要**认真处理共同资源**、**面对深层恐惧**，可能面临**财务压力**或需要**处理债务**。这是**建立资源管理能力**、**处理继承**、**学习情感深度**的好时期。适合**还债**、**处理保险**、**进行深层心理工作**。",
		"saturn_house_9": "**信念责任**和**哲学纪律**成为主题。你需要**认真对待信仰**、**建立哲学体系**，可能感到**信念受限**或需要**更深入地学习**。这是**完成高等教育**、**建立信仰体系**、**处理法律事务**的好时期。适合**深入学习**、**处理法律文件**、**建立精神传统**。",
		"saturn_house_10": "事业建构的关键时期。你会面对职业上的重大责任和挑战，可能需要付出更多努力。这是建立职业声誉、实现长期职业目标的重要阶段。成功需要耐心、纪律和持之以恒的努力。",
		"saturn_house_11": "**友谊责任**和**社群纪律**被强调。你需要**认真对待友谊**、**在团体中承担责任**，可能感到**社交受限**或需要**更成熟地参与社群**。这是**建立长期友谊**、**在团体中建立地位**、**实现集体目标**的好时期。适合**参与有意义的团体**、**建立社交网络**、**实现长期愿景**。",
		"saturn_house_12": "**内在责任**和**灵性纪律**成为主题。你需要**面对内在恐惧**、**处理业力**，可能感到**需要独处**或需要**更深入地修行**。这是**内在疗愈**、**释放旧模式**、**建立灵性纪律**的好时期。适合**冥想**、**处理隐藏议题**、**建立内在结构**。",
		
		// Mercury through houses - 完整 12 宫位
		"mercury_house_1": "**思维表达**和**个人沟通**成为焦点。你会**更频繁地表达观点**、**展现个人智慧**，**沟通能力**和**学习能力**提升。这是**展现个人魅力**、**建立个人品牌**、**提升表达能力**的好时期。适合**公开演讲**、**写作**、**学习新技能**。注意**沟通的准确性**，避免**过于急躁或散漫**。",
		"mercury_house_2": "**财务沟通**和**价值思考**被强调。你会**思考财务问题**、**讨论金钱话题**，**理财观念**和**价值判断**变得重要。这是**制定财务计划**、**学习理财知识**、**讨论合同条款**的好时期。适合**财务谈判**、**评估资产**、**学习投资**。注意**财务信息的准确性**，避免**冲动决策**。",
		"mercury_house_3": "**沟通交流**和**学习表达**达到高峰。你会**频繁交流**、**短途出行**，**学习**和**写作**都很顺利。这是**提升沟通技巧**、**建立人脉**、**学习新知识**的绝佳时期。适合**会议**、**写作**、**教学**、**短途旅行**。注意**信息的准确性**，避免**传播错误信息**。",
		"mercury_house_4": "**家庭沟通**和**情感表达**成为主题。你会**与家人更多交流**、**讨论家庭事务**，**家庭关系**和**情感表达**变得重要。这是**改善家庭沟通**、**处理家庭文件**、**学习家族历史**的好时期。适合**家庭会议**、**处理房产文件**、**记录家庭故事**。",
		"mercury_house_5": "**创意表达**和**娱乐沟通**被强调。你会**通过创作表达**、**在娱乐中交流**，**创造力**和**表达能力**结合。这是**写作**、**教学**、**创意项目**的好时期。适合**艺术创作**、**教育孩子**、**娱乐活动**。注意**表达的趣味性**，避免**过于严肃**。",
		"mercury_house_6": "**工作沟通**和**健康学习**成为重点。你会**在工作中频繁交流**、**学习健康知识**，**工作效率**和**健康管理**都受益。这是**改善工作流程**、**学习健康知识**、**处理工作文件**的好时期。适合**工作培训**、**健康咨询**、**优化工作方法**。",
		"mercury_house_7": "**关系沟通**和**合作交流**成为焦点。你会**与伙伴频繁交流**、**讨论合作关系**，**沟通技巧**在**关系中**变得重要。这是**改善关系沟通**、**签订合同**、**处理合作关系**的好时期。适合**关系咨询**、**商业谈判**、**建立合作**。注意**沟通的平衡**，避免**一方主导**。",
		"mercury_house_8": "**深度沟通**和**资源研究**被强调。你会**深入探讨**、**研究资源问题**，**心理洞察**和**资源分析**能力提升。这是**深度研究**、**处理共同资源**、**学习心理学**的好时期。适合**财务研究**、**心理咨询**、**处理遗产**。注意**信息的保密性**，避免**泄露敏感信息**。",
		"mercury_house_9": "**哲学思考**和**智慧传播**达到高峰。你会**思考人生意义**、**学习哲学**，**高等教育**和**远行**都很顺利。这是**学习**、**教学**、**旅行**、**法律事务**的绝佳时期。适合**高等教育**、**出版**、**法律咨询**、**跨文化学习**。",
		"mercury_house_10": "**职业沟通**和**公众表达**成为焦点。你会**在职业中频繁交流**、**展现专业能力**，**职业声誉**和**公众形象**受益。这是**职业演讲**、**专业写作**、**建立职业网络**的好时期。适合**职业发展**、**公众演讲**、**专业认证**。注意**职业形象的维护**。",
		"mercury_house_11": "**社群交流**和**理想沟通**被强调。你会**在团体中活跃交流**、**讨论理想愿景**，**社交网络**和**集体目标**受益。这是**参与团体活动**、**建立社交网络**、**传播理念**的好时期。适合**网络社交**、**团体讨论**、**实现集体愿景**。",
		"mercury_house_12": "**内在思考**和**潜意识沟通**成为主题。你会**深入思考**、**探索潜意识**，**直觉**和**灵感**变得活跃。这是**内在学习**、**灵性探索**、**处理隐藏议题**的好时期。适合**冥想**、**写作**、**心理疗愈**。注意**信息的清晰度**，避免**思维混乱**。",
		
		// Venus through houses - 补充剩余 9 个宫位
		"venus_house_1": "**个人魅力**和**审美提升**成为焦点。你会**更注重外表**、**展现个人魅力**，**吸引力**和**人缘**提升。这是**改善形象**、**提升魅力**、**享受生活**的好时期。适合**美容**、**购物**、**社交**。注意**保持真实**，避免**过度关注外表**。",
		"venus_house_2": "**财务享受**和**物质美感**被强调。你会**享受物质生活**、**投资于美的事物**，**收入**和**消费**都受益。这是**增加收入**、**享受物质**、**投资于价值**的好时期。适合**理财**、**购物**、**投资艺术品**。注意**理性消费**，避免**过度挥霍**。",
		"venus_house_3": "**沟通魅力**和**社交愉悦**成为主题。你会**在交流中展现魅力**、**享受社交**，**人际关系**和**学习**都受益。这是**改善沟通**、**建立人脉**、**享受学习**的好时期。适合**社交**、**学习**、**短途旅行**。",
		"venus_house_4": "**家庭和谐**和**情感美感**被强调。你会**美化家居**、**享受家庭时光**，**家庭关系**和**情感安全**受益。这是**改善家居**、**家庭聚会**、**享受家庭生活**的好时期。适合**家居装饰**、**家庭活动**、**处理房产**。",
		"venus_house_5": "浪漫和创造力达到高峰。这段时间充满魅力和吸引力，爱情关系可能开始或升温。你会更加享受生活的乐趣，艺术创作和娱乐活动带来愉悦。与孩子的互动也充满欢乐。这是追求浪漫、表达爱意、享受创造性活动的美好时期。",
		"venus_house_6": "**工作美感**和**健康享受**成为重点。你会**美化工作环境**、**享受健康生活**，**工作关系**和**身体健康**受益。这是**改善工作环境**、**享受工作**、**健康美容**的好时期。适合**工作社交**、**健康管理**、**享受日常**。",
		"venus_house_7": "关系和谐成为焦点。你会更加重视伙伴关系，渴望建立平衡和美好的连接。这是改善亲密关系、建立合作伙伴关系的好时机。你的外交技巧和魅力帮助你在人际关系中如鱼得水。适合解决关系冲突、签订合约。",
		"venus_house_8": "**深层吸引**和**资源美感**被强调。你会**在关系中深入**、**享受共同资源**，**亲密关系**和**财务共享**受益。这是**深化关系**、**处理共同财产**、**享受资源**的好时期。适合**关系深化**、**投资**、**处理遗产**。注意**关系的深度**，避免**过度依赖**。",
		"venus_house_9": "**精神美感**和**理想享受**成为主题。你会**享受精神追求**、**学习美学**，**信仰**和**旅行**都受益。这是**学习艺术**、**精神旅行**、**享受文化**的好时期。适合**艺术学习**、**文化旅行**、**精神探索**。",
		"venus_house_10": "**职业魅力**和**公众美感**成为焦点。你会**在职业中展现魅力**、**享受职业成就**，**职业关系**和**公众形象**受益。这是**提升职业形象**、**享受职业成功**、**建立职业关系**的好时期。适合**职业社交**、**公众形象**、**职业发展**。",
		"venus_house_11": "社交魅力大绽放，友谊和社群活动带来愉悦。你会结交新朋友，参与有趣的团体活动。社交网络扩大，可能通过朋友圈遇到浪漫对象。这段时间适合参与社交活动、追求集体目标、享受友谊的美好。",
		"venus_house_12": "**内在美感**和**灵性享受**被强调。你会**享受独处**、**探索内在美**，**灵性**和**艺术**受益。这是**内在美容**、**灵性艺术**、**享受宁静**的好时期。适合**冥想**、**艺术创作**、**内在探索**。",
		
		// Mars through houses - 补充剩余 9 个宫位
		"mars_house_1": "行动力和个人主动性大大增强。你会感到精力充沛，想要开始新的项目或挑战。自信心提升，但也可能变得冲动。这是实现个人目标、展现领导力的好时机，但要注意控制脾气，避免与他人冲突。",
		"mars_house_2": "**财务行动**和**物质竞争**被强调。你会**积极争取收入**、**竞争资源**，**赚钱能力**和**消费冲动**都增强。这是**增加收入**、**投资**、**争取财务目标**的好时期。适合**财务行动**、**投资**、**争取加薪**。注意**理性消费**，避免**冲动购买**。",
		"mars_house_3": "**沟通行动**和**学习竞争**成为主题。你会**积极交流**、**竞争学习**，**沟通能力**和**学习动力**提升。这是**积极沟通**、**学习竞争**、**短途出行**的好时期。适合**辩论**、**学习**、**短途旅行**。注意**沟通的语气**，避免**过于激进**。",
		"mars_house_4": "**家庭行动**和**情感竞争**被强调。你会**积极处理家庭事务**、**争取家庭地位**，**家庭关系**和**情感表达**变得激烈。这是**改善家庭**、**处理家庭冲突**、**争取家庭目标**的好时期。适合**家庭项目**、**处理房产**、**家庭活动**。注意**家庭和谐**，避免**家庭冲突**。",
		"mars_house_5": "**创意行动**和**娱乐竞争**成为焦点。你会**积极创作**、**竞争娱乐**，**创造力**和**行动力**结合。这是**积极创作**、**运动**、**竞争娱乐**的好时期。适合**艺术创作**、**运动**、**娱乐竞争**。注意**安全**，避免**过度冒险**。",
		"mars_house_6": "**工作行动**和**健康竞争**成为重点。你会**积极工作**、**竞争健康**，**工作效率**和**健康管理**都受益。这是**积极工作**、**运动健身**、**竞争工作目标**的好时期。适合**工作冲刺**、**健身**、**改善健康**。注意**工作平衡**，避免**过度劳累**。",
		"mars_house_7": "**关系行动**和**合作竞争**被强调。你会**积极处理关系**、**竞争合作关系**，**关系动力**和**合作能力**提升。这是**积极改善关系**、**竞争合作**、**处理关系冲突**的好时期。适合**关系行动**、**合作竞争**、**处理关系**。注意**关系平衡**，避免**关系冲突**。",
		"mars_house_8": "**深度行动**和**资源竞争**成为主题。你会**积极处理资源**、**竞争深度关系**，**资源管理**和**情感深度**变得激烈。这是**积极处理资源**、**竞争投资**、**处理共同财产**的好时期。适合**资源行动**、**投资竞争**、**处理债务**。注意**资源安全**，避免**过度冒险**。",
		"mars_house_9": "**精神行动**和**理想竞争**被强调。你会**积极追求理想**、**竞争精神目标**，**精神追求**和**理想实现**受益。这是**积极学习**、**竞争精神**、**追求理想**的好时期。适合**学习竞争**、**精神追求**、**理想行动**。",
		"mars_house_10": "事业野心和职业推动力达到高峰。你会全力以赴追求职业目标，可能会有重要的职业突破。工作中展现出强大的执行力和竞争力。这是争取晋升、推进重要项目的好时机，但要注意与权威人物的关系。",
		"mars_house_11": "积极勇敢地追求理想和社会目标。你会参与团体活动，为共同的目标而奋斗。社交中展现出行动力和领导能力。这段时间适合推动社会变革、实现集体目标，但要注意团队合作，避免独断专行。",
		"mars_house_12": "**内在行动**和**灵性竞争**成为主题。你会**积极处理内在议题**、**竞争灵性目标**，**内在动力**和**灵性追求**受益。这是**积极疗愈**、**内在行动**、**灵性竞争**的好时期。适合**内在工作**、**灵性行动**、**处理隐藏议题**。注意**内在平衡**，避免**内在冲突**。",
		
		// Jupiter through houses - 补充剩余 9 个宫位
		"jupiter_house_1": "**个人扩展**和**乐观自信**成为焦点。你会**感到自信**、**展现个人魅力**，**个人成长**和**机会**增多。这是**个人发展**、**展现魅力**、**抓住机会**的好时期。适合**个人目标**、**展现能力**、**抓住机遇**。注意**保持谦逊**，避免**过度自信**。",
		"jupiter_house_2": "财富与机遇增加的时期。收入可能提升，或是获得新的赚钱机会。你对物质的态度变得更加乐观和慷慨。这是投资、扩展财务资源的好时机。注意不要过度消费或过于乐观，保持理性的财务规划。",
		"jupiter_house_3": "**沟通扩展**和**学习机会**被强调。你会**频繁交流**、**学习机会**增多，**沟通能力**和**学习能力**提升。这是**学习**、**交流**、**短途旅行**的好时期。适合**学习**、**教学**、**建立人脉**。",
		"jupiter_house_4": "**家庭扩展**和**情感乐观**成为主题。你会**家庭关系**改善、**家庭机会**增多，**家庭关系**和**情感安全**受益。这是**改善家庭**、**家庭扩展**、**情感乐观**的好时期。适合**家庭活动**、**处理房产**、**家庭成长**。",
		"jupiter_house_5": "**创意扩展**和**娱乐机会**被强调。你会**创作机会**增多、**娱乐享受**提升，**创造力**和**娱乐能力**受益。这是**创作**、**娱乐**、**教育孩子**的好时期。适合**艺术创作**、**娱乐活动**、**教育**。",
		"jupiter_house_6": "**工作扩展**和**健康机会**成为重点。你会**工作机会**增多、**健康改善**，**工作效率**和**健康管理**受益。这是**工作发展**、**健康改善**、**工作机会**的好时期。适合**工作扩展**、**健康管理**、**工作学习**。",
		"jupiter_house_7": "**关系扩展**和**合作机会**成为焦点。你会**关系机会**增多、**合作扩展**，**关系质量**和**合作能力**受益。这是**关系发展**、**合作扩展**、**关系机会**的好时期。适合**关系发展**、**合作**、**关系学习**。",
		"jupiter_house_8": "**资源扩展**和**深度机会**被强调。你会**资源机会**增多、**深度关系**扩展，**资源管理**和**情感深度**受益。这是**资源扩展**、**投资机会**、**深度关系**的好时期。适合**投资**、**资源管理**、**深度关系**。",
		"jupiter_house_9": "智慧启蒙和精神成长的黄金时期。你会对哲学、宗教、高等教育产生浓厚兴趣。可能有机会远行或接触不同文化。视野开阔，对生命意义有更深的理解。这是学习、旅行、追求真理的绝佳时机。",
		"jupiter_house_10": "**职业扩展**和**事业机会**成为焦点。你会**职业机会**增多、**事业扩展**，**职业发展**和**事业成就**受益。这是**职业发展**、**事业扩展**、**职业机会**的好时期。适合**职业发展**、**事业扩展**、**职业学习**。",
		"jupiter_house_11": "贵人现身还请抓住机会。友谊和社交活动带来幸运和机遇。你会结识有影响力的朋友，社交网络扩大。参与团体活动能带来成长和收获。这是实现长期目标、获得社会支持的好时期。",
		"jupiter_house_12": "**灵性扩展**和**内在机会**被强调。你会**灵性机会**增多、**内在成长**扩展，**灵性追求**和**内在探索**受益。这是**灵性发展**、**内在扩展**、**灵性机会**的好时期。适合**灵性学习**、**内在探索**、**灵性成长**。",
		
		// Uranus through houses - 完整 12 宫位
		"uranus_house_1": "**个人突破**和**突然改变**成为焦点。你会**突然改变形象**、**突破个人限制**，**个人自由**和**独立性**增强。这是**个人突破**、**改变形象**、**追求自由**的好时期。适合**个人改变**、**突破限制**、**追求独立**。注意**变化的稳定性**，避免**过于突然**。",
		"uranus_house_2": "**财务突破**和**价值改变**被强调。你会**突然改变财务**、**突破价值观念**，**财务自由**和**价值独立**增强。这是**财务突破**、**改变价值**、**追求财务自由**的好时期。适合**财务创新**、**价值改变**、**财务独立**。注意**财务安全**，避免**过于冒险**。",
		"uranus_house_3": "**沟通突破**和**思维改变**成为主题。你会**突然改变思维**、**突破沟通方式**，**思维自由**和**沟通独立**增强。这是**思维突破**、**改变沟通**、**追求思维自由**的好时期。适合**思维创新**、**沟通改变**、**学习新技术**。",
		"uranus_house_4": "**家庭突破**和**情感改变**被强调。你会**突然改变家庭**、**突破家庭模式**，**家庭自由**和**情感独立**增强。这是**家庭突破**、**改变家庭**、**追求家庭自由**的好时期。适合**家庭改变**、**处理房产**、**家庭独立**。注意**家庭稳定**，避免**过于突然**。",
		"uranus_house_5": "**创意突破**和**娱乐改变**成为焦点。你会**突然改变创作**、**突破娱乐方式**，**创意自由**和**娱乐独立**增强。这是**创意突破**、**改变娱乐**、**追求创意自由**的好时期。适合**创意创新**、**娱乐改变**、**创意独立**。",
		"uranus_house_6": "**工作突破**和**健康改变**成为重点。你会**突然改变工作**、**突破工作方式**，**工作自由**和**健康独立**增强。这是**工作突破**、**改变工作**、**追求工作自由**的好时期。适合**工作创新**、**健康改变**、**工作独立**。注意**工作稳定**，避免**过于突然**。",
		"uranus_house_7": "**关系突破**和**合作改变**被强调。你会**突然改变关系**、**突破合作方式**，**关系自由**和**合作独立**增强。这是**关系突破**、**改变关系**、**追求关系自由**的好时期。适合**关系创新**、**合作改变**、**关系独立**。注意**关系稳定**，避免**过于突然**。",
		"uranus_house_8": "**资源突破**和**深度改变**成为主题。你会**突然改变资源**、**突破深度关系**，**资源自由**和**深度独立**增强。这是**资源突破**、**改变资源**、**追求资源自由**的好时期。适合**资源创新**、**投资改变**、**资源独立**。注意**资源安全**，避免**过于冒险**。",
		"uranus_house_9": "**精神突破**和**理想改变**被强调。你会**突然改变理想**、**突破精神方式**，**精神自由**和**理想独立**增强。这是**精神突破**、**改变理想**、**追求精神自由**的好时期。适合**精神创新**、**理想改变**、**精神独立**。",
		"uranus_house_10": "**职业突破**和**事业改变**成为焦点。你会**突然改变职业**、**突破事业方式**，**职业自由**和**事业独立**增强。这是**职业突破**、**改变职业**、**追求职业自由**的好时期。适合**职业创新**、**事业改变**、**职业独立**。注意**职业稳定**，避免**过于突然**。",
		"uranus_house_11": "**社群突破**和**理想改变**被强调。你会**突然改变社群**、**突破理想方式**，**社群自由**和**理想独立**增强。这是**社群突破**、**改变社群**、**追求社群自由**的好时期。适合**社群创新**、**理想改变**、**社群独立**。",
		"uranus_house_12": "**内在突破**和**灵性改变**成为主题。你会**突然改变内在**、**突破灵性方式**，**内在自由**和**灵性独立**增强。这是**内在突破**、**改变内在**、**追求内在自由**的好时期。适合**内在创新**、**灵性改变**、**内在独立**。",
		
		// Neptune through houses - 完整 12 宫位
		"neptune_house_1": "**个人理想**和**灵性表达**成为焦点。你会**展现灵性魅力**、**追求理想自我**，**个人灵感**和**灵性表达**增强。这是**个人理想**、**灵性表达**、**追求理想**的好时期。适合**灵性探索**、**理想表达**、**个人灵感**。注意**保持现实**，避免**过度理想化**。",
		"neptune_house_2": "**财务理想**和**价值灵感**被强调。你会**理想化财务**、**追求价值灵感**，**财务理想**和**价值灵感**增强。这是**财务理想**、**价值灵感**、**追求理想价值**的好时期。适合**理想投资**、**价值灵感**、**财务理想**。注意**财务现实**，避免**过度理想化**。",
		"neptune_house_3": "**沟通理想**和**思维灵感**成为主题。你会**理想化沟通**、**追求思维灵感**，**沟通理想**和**思维灵感**增强。这是**沟通理想**、**思维灵感**、**追求理想沟通**的好时期。适合**理想沟通**、**思维灵感**、**沟通理想**。注意**沟通现实**，避免**过度理想化**。",
		"neptune_house_4": "**家庭理想**和**情感灵感**被强调。你会**理想化家庭**、**追求情感灵感**，**家庭理想**和**情感灵感**增强。这是**家庭理想**、**情感灵感**、**追求理想家庭**的好时期。适合**理想家庭**、**情感灵感**、**家庭理想**。注意**家庭现实**，避免**过度理想化**。",
		"neptune_house_5": "**创意理想**和**娱乐灵感**成为焦点。你会**理想化创作**、**追求娱乐灵感**，**创意理想**和**娱乐灵感**增强。这是**创意理想**、**娱乐灵感**、**追求理想创作**的好时期。适合**理想创作**、**娱乐灵感**、**创意理想**。",
		"neptune_house_6": "**工作理想**和**健康灵感**成为重点。你会**理想化工作**、**追求健康灵感**，**工作理想**和**健康灵感**增强。这是**工作理想**、**健康灵感**、**追求理想工作**的好时期。适合**理想工作**、**健康灵感**、**工作理想**。注意**工作现实**，避免**过度理想化**。",
		"neptune_house_7": "**关系理想**和**合作灵感**被强调。你会**理想化关系**、**追求合作灵感**，**关系理想**和**合作灵感**增强。这是**关系理想**、**合作灵感**、**追求理想关系**的好时期。适合**理想关系**、**合作灵感**、**关系理想**。注意**关系现实**，避免**过度理想化**。",
		"neptune_house_8": "**资源理想**和**深度灵感**成为主题。你会**理想化资源**、**追求深度灵感**，**资源理想**和**深度灵感**增强。这是**资源理想**、**深度灵感**、**追求理想资源**的好时期。适合**理想资源**、**深度灵感**、**资源理想**。注意**资源现实**，避免**过度理想化**。",
		"neptune_house_9": "**精神理想**和**哲学灵感**被强调。你会**理想化精神**、**追求哲学灵感**，**精神理想**和**哲学灵感**增强。这是**精神理想**、**哲学灵感**、**追求理想精神**的好时期。适合**理想精神**、**哲学灵感**、**精神理想**。",
		"neptune_house_10": "**职业理想**和**事业灵感**成为焦点。你会**理想化职业**、**追求事业灵感**，**职业理想**和**事业灵感**增强。这是**职业理想**、**事业灵感**、**追求理想职业**的好时期。适合**理想职业**、**事业灵感**、**职业理想**。注意**职业现实**，避免**过度理想化**。",
		"neptune_house_11": "**社群理想**和**集体灵感**被强调。你会**理想化社群**、**追求集体灵感**，**社群理想**和**集体灵感**增强。这是**社群理想**、**集体灵感**、**追求理想社群**的好时期。适合**理想社群**、**集体灵感**、**社群理想**。",
		"neptune_house_12": "**内在理想**和**灵性灵感**成为主题。你会**理想化内在**、**追求灵性灵感**，**内在理想**和**灵性灵感**增强。这是**内在理想**、**灵性灵感**、**追求理想内在**的好时期。适合**理想内在**、**灵性灵感**、**内在理想**。",
		
		// Pluto through houses - 完整 12 宫位
		"pluto_house_1": "**个人转化**和**深层改变**成为焦点。你会**经历深层转化**、**改变个人本质**，**个人力量**和**转化能力**增强。这是**个人转化**、**深层改变**、**追求个人力量**的好时期。适合**个人转化**、**深层改变**、**个人力量**。注意**转化的稳定性**，避免**过于激烈**。",
		"pluto_house_2": "**财务转化**和**价值改变**被强调。你会**经历财务转化**、**改变价值观念**，**财务力量**和**价值转化**增强。这是**财务转化**、**价值改变**、**追求财务力量**的好时期。适合**财务转化**、**价值改变**、**财务力量**。注意**财务安全**，避免**过于激烈**。",
		"pluto_house_3": "**沟通转化**和**思维改变**成为主题。你会**经历沟通转化**、**改变思维方式**，**沟通力量**和**思维转化**增强。这是**沟通转化**、**思维改变**、**追求沟通力量**的好时期。适合**沟通转化**、**思维改变**、**沟通力量**。",
		"pluto_house_4": "**家庭转化**和**情感改变**被强调。你会**经历家庭转化**、**改变家庭模式**，**家庭力量**和**情感转化**增强。这是**家庭转化**、**情感改变**、**追求家庭力量**的好时期。适合**家庭转化**、**情感改变**、**家庭力量**。注意**家庭稳定**，避免**过于激烈**。",
		"pluto_house_5": "**创意转化**和**娱乐改变**成为焦点。你会**经历创意转化**、**改变娱乐方式**，**创意力量**和**娱乐转化**增强。这是**创意转化**、**娱乐改变**、**追求创意力量**的好时期。适合**创意转化**、**娱乐改变**、**创意力量**。",
		"pluto_house_6": "**工作转化**和**健康改变**成为重点。你会**经历工作转化**、**改变工作方式**，**工作力量**和**健康转化**增强。这是**工作转化**、**健康改变**、**追求工作力量**的好时期。适合**工作转化**、**健康改变**、**工作力量**。注意**工作稳定**，避免**过于激烈**。",
		"pluto_house_7": "**关系转化**和**合作改变**被强调。你会**经历关系转化**、**改变合作方式**，**关系力量**和**合作转化**增强。这是**关系转化**、**合作改变**、**追求关系力量**的好时期。适合**关系转化**、**合作改变**、**关系力量**。注意**关系稳定**，避免**过于激烈**。",
		"pluto_house_8": "**资源转化**和**深度改变**成为主题。你会**经历资源转化**、**改变深度关系**，**资源力量**和**深度转化**增强。这是**资源转化**、**深度改变**、**追求资源力量**的好时期。适合**资源转化**、**投资改变**、**资源力量**。注意**资源安全**，避免**过于激烈**。",
		"pluto_house_9": "**精神转化**和**理想改变**被强调。你会**经历精神转化**、**改变理想方式**，**精神力量**和**理想转化**增强。这是**精神转化**、**理想改变**、**追求精神力量**的好时期。适合**精神转化**、**理想改变**、**精神力量**。",
		"pluto_house_10": "**职业转化**和**事业改变**成为焦点。你会**经历职业转化**、**改变事业方式**，**职业力量**和**事业转化**增强。这是**职业转化**、**事业改变**、**追求职业力量**的好时期。适合**职业转化**、**事业改变**、**职业力量**。注意**职业稳定**，避免**过于激烈**。",
		"pluto_house_11": "**社群转化**和**理想改变**被强调。你会**经历社群转化**、**改变理想方式**，**社群力量**和**理想转化**增强。这是**社群转化**、**理想改变**、**追求社群力量**的好时期。适合**社群转化**、**理想改变**、**社群力量**。",
		"pluto_house_12": "**内在转化**和**灵性改变**成为主题。你会**经历内在转化**、**改变灵性方式**，**内在力量**和**灵性转化**增强。这是**内在转化**、**灵性改变**、**追求内在力量**的好时期。适合**内在转化**、**灵性改变**、**内在力量**。",
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
		"moon_house_2": "Emotional security is closely tied to material security. You'll value stability and comfort more, possibly satisfying emotional needs through shopping or food. Good time to organize finances, assess resources, or invest in things that bring security. You're particularly sensitive to beauty and can enjoy sensory pleasures. Avoid overspending to fill emotional voids.",
		"moon_house_3": "Emotional expression and communication become active. You'll have more need to communicate, wanting to share feelings and thoughts. Interactions with siblings and neighbors increase; short trips are frequent. Good time for journaling, reading, learning. Your thinking is influenced by emotions, possibly viewing things more emotionally. Maintain sincerity and empathy in communication.",
		"moon_house_4": "Home belonging and emotional roots become the focus. You'll crave familiar environments, with deeper emotional connections to family. Childhood memories and family patterns may surface, affecting current emotions. Good time to care for family, organize home, or handle family relationships. You need emotional nourishment and security.",
		"moon_house_5": "Emotional expression becomes passionate and creative. You crave romance, entertainment, and creative activities. Interaction with children or personal hobbies brings emotional satisfaction. Good time to pursue joy, express love, engage in artistic creation. Your emotions may be dramatic but full of vitality and passion.",
		"moon_house_6": "Emotions affect physical health and daily work. You'll pay more attention to health habits, possibly stabilizing emotions through regular routines. Work needs emotional support; relationships with colleagues affect work state. Good time to adjust lifestyle, focus on health, establish daily routines. Serving others brings emotional satisfaction.",
		"moon_house_7": "Emotional needs are met through one-on-one relationships. You'll depend more on intimate partners, craving emotional resonance and understanding. Emotional dynamics in relationships become obvious; you may be more easily affected by others' emotions. Good time to deepen intimate relationships, seek partner support, but maintain emotional independence.",
		"moon_house_8": "Emotions go deep; you'll face intense emotional experiences. May experience emotional transformation or deep psychological purification. Emotional bonds with others deepen; shared resources and intimacy increase. Good time to explore inner depths, face fears, engage in emotional healing. Emotions may be intense; need safe expression.",
		"moon_house_9": "Emotions expand to broader fields. You crave learning and travel to satisfy emotional needs, developing interest in different cultures and philosophies. Emotions become optimistic and open but may be unstable. Good time to explore new beliefs, expand horizons, nourish emotions through spiritual growth.",
		"moon_house_10": "Public image and emotional needs intersect. Your emotions may show in public; career is affected by emotions. You crave recognition at work to satisfy emotional needs. Good time to show your authentic self in career, but also mind professional image management.",
		"moon_house_11": "Friendships and community activities satisfy emotional needs. You'll participate more in group activities, getting emotional support from friends. Future visions bring emotional inspiration. Good time to develop friendships, participate in community, pursue collective goals. Emotionally need belonging and shared ideals.",
		"moon_house_12": "Emotions enter subconscious level; you may feel emotionally vague and sensitive. Need solitude to process inner emotions; dreams and intuition become active. Good time for meditation, solitude, emotional healing. You may feel emotionally vulnerable; need space and time to recover.",
		
		// Mercury through houses - complete 12 houses
		"mercury_house_1": "**Mental expression** and **personal communication** are highlighted. You'll **express views more frequently**, **show personal wisdom**; **communication** and **learning abilities** improve. Good time to **show personal charm**, **build personal brand**, **enhance expression**. Suitable for **public speaking**, **writing**, **learning new skills**. Mind **communication accuracy**; avoid **hastiness or scatter**.",
		"mercury_house_2": "**Financial communication** and **value thinking** are emphasized. You'll **think about finances**, **discuss money**; **financial concepts** and **value judgment** become important. Good time for **financial planning**, **learning finance**, **discussing contracts**. Suitable for **financial negotiation**, **asset assessment**, **learning investment**. Mind **financial accuracy**; avoid **impulsive decisions**.",
		"mercury_house_3": "**Communication** and **learning expression** peak. You'll **communicate frequently**, **take short trips**; **learning** and **writing** flow smoothly. Excellent time to **improve communication**, **build networks**, **learn new knowledge**. Suitable for **meetings**, **writing**, **teaching**, **short travel**. Mind **information accuracy**; avoid **spreading misinformation**.",
		"mercury_house_4": "**Family communication** and **emotional expression** are themes. You'll **communicate more with family**, **discuss family matters**; **family relationships** and **emotional expression** become important. Good time to **improve family communication**, **handle family documents**, **learn family history**. Suitable for **family meetings**, **property documents**, **recording family stories**.",
		"mercury_house_5": "**Creative expression** and **entertainment communication** are emphasized. You'll **express through creation**, **communicate in entertainment**; **creativity** and **expression** combine. Good time for **writing**, **teaching**, **creative projects**. Suitable for **artistic creation**, **educating children**, **entertainment**. Mind **expression interest**; avoid **over-seriousness**.",
		"mercury_house_6": "**Work communication** and **health learning** are priorities. You'll **communicate frequently at work**, **learn health knowledge**; **work efficiency** and **health management** benefit. Good time to **improve work processes**, **learn health**, **handle work documents**. Suitable for **work training**, **health consultation**, **optimizing work methods**.",
		"mercury_house_7": "**Relationship communication** and **cooperation exchange** are highlighted. You'll **communicate frequently with partners**, **discuss cooperation**; **communication skills** in **relationships** become important. Good time to **improve relationship communication**, **sign contracts**, **handle cooperation**. Suitable for **relationship consultation**, **business negotiation**, **building cooperation**. Mind **communication balance**; avoid **one-sided dominance**.",
		"mercury_house_8": "**Deep communication** and **resource research** are emphasized. You'll **explore deeply**, **research resources**; **psychological insight** and **resource analysis** improve. Good time for **deep research**, **handling shared resources**, **learning psychology**. Suitable for **financial research**, **psychological consultation**, **handling inheritance**. Mind **information confidentiality**; avoid **leaking sensitive info**.",
		"mercury_house_9": "**Philosophical thinking** and **wisdom spreading** peak. You'll **think about life meaning**, **study philosophy**; **higher education** and **travel** flow smoothly. Excellent time for **learning**, **teaching**, **travel**, **legal matters**. Suitable for **higher education**, **publishing**, **legal consultation**, **cross-cultural learning**.",
		"mercury_house_10": "**Career communication** and **public expression** are highlighted. You'll **communicate frequently in career**, **show professional ability**; **career reputation** and **public image** benefit. Good time for **career speeches**, **professional writing**, **building career networks**. Suitable for **career development**, **public speaking**, **professional certification**. Mind **career image maintenance**.",
		"mercury_house_11": "**Community exchange** and **ideal communication** are emphasized. You'll **communicate actively in groups**, **discuss ideal visions**; **social networks** and **collective goals** benefit. Good time to **participate in group activities**, **build social networks**, **spread ideas**. Suitable for **online socializing**, **group discussions**, **realizing collective visions**.",
		"mercury_house_12": "**Inner thinking** and **subconscious communication** are themes. You'll **think deeply**, **explore subconscious**; **intuition** and **inspiration** become active. Good time for **inner learning**, **spiritual exploration**, **handling hidden issues**. Suitable for **meditation**, **writing**, **psychological healing**. Mind **information clarity**; avoid **mental confusion**.",
		
		// Venus through houses - complete remaining 9 houses
		"venus_house_1": "**Personal charm** and **aesthetic improvement** are highlighted. You'll **focus more on appearance**, **show personal charm**; **attractiveness** and **rapport** improve. Good time to **improve image**, **enhance charm**, **enjoy life**. Suitable for **beauty**, **shopping**, **socializing**. Mind **staying authentic**; avoid **over-focusing on appearance**.",
		"venus_house_2": "**Financial enjoyment** and **material beauty** are emphasized. You'll **enjoy material life**, **invest in beautiful things**; **income** and **spending** benefit. Good time to **increase income**, **enjoy material**, **invest in value**. Suitable for **finance**, **shopping**, **investing in art**. Mind **rational spending**; avoid **overspending**.",
		"venus_house_3": "**Communication charm** and **social pleasure** are themes. You'll **show charm in communication**, **enjoy socializing**; **relationships** and **learning** benefit. Good time to **improve communication**, **build networks**, **enjoy learning**. Suitable for **socializing**, **learning**, **short travel**.",
		"venus_house_4": "**Family harmony** and **emotional beauty** are emphasized. You'll **beautify home**, **enjoy family time**; **family relationships** and **emotional security** benefit. Good time to **improve home**, **family gatherings**, **enjoy family life**. Suitable for **home decoration**, **family activities**, **handling property**.",
		"venus_house_5": "Romance and creativity peak. This period is full of charm and attraction; romantic relationships may begin or heat up. You'll enjoy life's pleasures more; artistic creation and entertainment bring joy. Interaction with children is also joyful. Beautiful time to pursue romance, express love, enjoy creative activities.",
		"venus_house_6": "**Work beauty** and **health enjoyment** are priorities. You'll **beautify work environment**, **enjoy healthy life**; **work relationships** and **physical health** benefit. Good time to **improve work environment**, **enjoy work**, **health and beauty**. Suitable for **work socializing**, **health management**, **enjoying daily life**.",
		"venus_house_7": "Relationship harmony becomes the focus. You'll value partnerships more, craving balanced and beautiful connections. Good time to improve intimate relationships, establish partnerships. Your diplomatic skills and charm help you navigate relationships smoothly. Suitable for resolving relationship conflicts, signing contracts.",
		"venus_house_8": "**Deep attraction** and **resource beauty** are emphasized. You'll **deepen in relationships**, **enjoy shared resources**; **intimate relationships** and **financial sharing** benefit. Good time to **deepen relationships**, **handle shared property**, **enjoy resources**. Suitable for **relationship deepening**, **investment**, **handling inheritance**. Mind **relationship depth**; avoid **over-dependence**.",
		"venus_house_9": "**Spiritual beauty** and **ideal enjoyment** are themes. You'll **enjoy spiritual pursuits**, **study aesthetics**; **faith** and **travel** benefit. Good time to **study art**, **spiritual travel**, **enjoy culture**. Suitable for **artistic learning**, **cultural travel**, **spiritual exploration**.",
		"venus_house_10": "**Career charm** and **public beauty** are highlighted. You'll **show charm in career**, **enjoy career achievement**; **career relationships** and **public image** benefit. Good time to **enhance career image**, **enjoy career success**, **build career relationships**. Suitable for **career socializing**, **public image**, **career development**.",
		"venus_house_11": "Social charisma blooms as friendships and community activities bring joy. You'll make new friends and participate in interesting group activities. Your social network expands, and you might meet romantic interests through friends. This is a great time to engage in social activities, pursue collective goals, and enjoy the beauty of friendship.",
		"venus_house_12": "**Inner beauty** and **spiritual enjoyment** are emphasized. You'll **enjoy solitude**, **explore inner beauty**; **spirituality** and **arts** benefit. Good time for **inner beauty**, **spiritual art**, **enjoying peace**. Suitable for **meditation**, **artistic creation**, **inner exploration**.",
		
		// Mars through houses - complete remaining 9 houses
		"mars_house_1": "Drive and personal initiative greatly increase. You'll feel energetic, wanting to start new projects or challenges. Confidence rises, but may become impulsive. Good time to achieve personal goals, show leadership, but mind temper control, avoid conflicts.",
		"mars_house_2": "**Financial action** and **material competition** are emphasized. You'll **actively pursue income**, **compete for resources**; **earning ability** and **spending impulse** increase. Good time to **increase income**, **invest**, **pursue financial goals**. Suitable for **financial action**, **investment**, **seeking raises**. Mind **rational spending**; avoid **impulse buying**.",
		"mars_house_3": "**Communication action** and **learning competition** are themes. You'll **communicate actively**, **compete in learning**; **communication** and **learning drive** improve. Good time for **active communication**, **learning competition**, **short travel**. Suitable for **debate**, **learning**, **short trips**. Mind **communication tone**; avoid **over-aggression**.",
		"mars_house_4": "**Family action** and **emotional competition** are emphasized. You'll **handle family matters actively**, **compete for family status**; **family relationships** and **emotional expression** become intense. Good time to **improve family**, **handle family conflicts**, **pursue family goals**. Suitable for **family projects**, **handling property**, **family activities**. Mind **family harmony**; avoid **family conflicts**.",
		"mars_house_5": "**Creative action** and **entertainment competition** are highlighted. You'll **create actively**, **compete in entertainment**; **creativity** and **action** combine. Good time for **active creation**, **sport**, **entertainment competition**. Suitable for **artistic creation**, **sport**, **entertainment competition**. Mind **safety**; avoid **excessive risk**.",
		"mars_house_6": "**Work action** and **health competition** are priorities. You'll **work actively**, **compete for health**; **work efficiency** and **health management** benefit. Good time for **active work**, **exercise**, **competing for work goals**. Suitable for **work sprints**, **fitness**, **health improvement**. Mind **work balance**; avoid **overwork**.",
		"mars_house_7": "**Relationship action** and **cooperation competition** are emphasized. You'll **handle relationships actively**, **compete for cooperation**; **relationship drive** and **cooperation ability** improve. Good time to **improve relationships actively**, **compete for cooperation**, **handle relationship conflicts**. Suitable for **relationship action**, **cooperation competition**, **handling relationships**. Mind **relationship balance**; avoid **relationship conflicts**.",
		"mars_house_8": "**Deep action** and **resource competition** are themes. You'll **handle resources actively**, **compete for deep relationships**; **resource management** and **emotional depth** become intense. Good time for **active resource handling**, **investment competition**, **handling shared property**. Suitable for **resource action**, **investment competition**, **handling debt**. Mind **resource safety**; avoid **excessive risk**.",
		"mars_house_9": "**Spiritual action** and **ideal competition** are emphasized. You'll **pursue ideals actively**, **compete for spiritual goals**; **spiritual pursuit** and **ideal realization** benefit. Good time for **active learning**, **spiritual competition**, **pursuing ideals**. Suitable for **learning competition**, **spiritual pursuit**, **ideal action**.",
		"mars_house_10": "Career ambition and professional drive peak. You'll go all out for career goals, possibly major career breakthroughs. Work shows strong execution and competitiveness. Good time to seek promotion, advance key projects, but mind relationships with authority.",
		"mars_house_11": "Bold and brave pursuit of ideals and social goals. You'll participate in group activities, fighting for common objectives. You show drive and leadership in social settings. This period is suitable for promoting social change and achieving collective goals, but remember to cooperate with the team.",
		"mars_house_12": "**Inner action** and **spiritual competition** are themes. You'll **handle inner issues actively**, **compete for spiritual goals**; **inner drive** and **spiritual pursuit** benefit. Good time for **active healing**, **inner action**, **spiritual competition**. Suitable for **inner work**, **spiritual action**, **handling hidden issues**. Mind **inner balance**; avoid **inner conflicts**.",
		
		// Jupiter through houses - complete remaining 9 houses
		"jupiter_house_1": "**Personal expansion** and **optimistic confidence** are highlighted. You'll **feel confident**, **show personal charm**; **personal growth** and **opportunities** increase. Good time for **personal development**, **showing charm**, **seizing opportunities**. Suitable for **personal goals**, **showing ability**, **seizing opportunities**. Mind **staying humble**; avoid **over-confidence**.",
		"jupiter_house_2": "Period of increased wealth and opportunities. Income may rise, or gain new earning chances. Your attitude toward material becomes more optimistic and generous. Good time to invest, expand financial resources. Mind not overspending or over-optimism; maintain rational financial planning.",
		"jupiter_house_3": "**Communication expansion** and **learning opportunities** are emphasized. You'll **communicate frequently**, **learning opportunities** increase; **communication** and **learning abilities** improve. Good time for **learning**, **exchange**, **short travel**. Suitable for **learning**, **teaching**, **building networks**.",
		"jupiter_house_4": "**Family expansion** and **emotional optimism** are themes. You'll **family relationships** improve, **family opportunities** increase; **family relationships** and **emotional security** benefit. Good time to **improve family**, **family expansion**, **emotional optimism**. Suitable for **family activities**, **handling property**, **family growth**.",
		"jupiter_house_5": "**Creative expansion** and **entertainment opportunities** are emphasized. You'll **creative opportunities** increase, **entertainment enjoyment** improve; **creativity** and **entertainment ability** benefit. Good time for **creation**, **entertainment**, **educating children**. Suitable for **artistic creation**, **entertainment activities**, **education**.",
		"jupiter_house_6": "**Work expansion** and **health opportunities** are priorities. You'll **work opportunities** increase, **health improvement**; **work efficiency** and **health management** benefit. Good time for **work development**, **health improvement**, **work opportunities**. Suitable for **work expansion**, **health management**, **work learning**.",
		"jupiter_house_7": "**Relationship expansion** and **cooperation opportunities** are highlighted. You'll **relationship opportunities** increase, **cooperation expansion**; **relationship quality** and **cooperation ability** benefit. Good time for **relationship development**, **cooperation expansion**, **relationship opportunities**. Suitable for **relationship development**, **cooperation**, **relationship learning**.",
		"jupiter_house_8": "**Resource expansion** and **deep opportunities** are emphasized. You'll **resource opportunities** increase, **deep relationship** expansion; **resource management** and **emotional depth** benefit. Good time for **resource expansion**, **investment opportunities**, **deep relationships**. Suitable for **investment**, **resource management**, **deep relationships**.",
		"jupiter_house_9": "Golden period of wisdom enlightenment and spiritual growth. You'll develop strong interest in philosophy, religion, higher education. May have chances to travel far or contact different cultures. Vision broadens; deeper understanding of life meaning. Excellent time for learning, travel, pursuing truth.",
		"jupiter_house_10": "**Career expansion** and **career opportunities** are highlighted. You'll **career opportunities** increase, **career expansion**; **career development** and **career achievement** benefit. Good time for **career development**, **career expansion**, **career opportunities**. Suitable for **career development**, **career expansion**, **career learning**.",
		"jupiter_house_11": "Benefactors appear - seize the opportunity. Friendships and social activities bring luck and opportunities. You'll meet influential friends, expanding your social network. Participating in group activities brings growth and rewards. This is a good period for achieving long-term goals and gaining social support.",
		"jupiter_house_12": "**Spiritual expansion** and **inner opportunities** are emphasized. You'll **spiritual opportunities** increase, **inner growth** expansion; **spiritual pursuit** and **inner exploration** benefit. Good time for **spiritual development**, **inner expansion**, **spiritual opportunities**. Suitable for **spiritual learning**, **inner exploration**, **spiritual growth**.",
		
		// Saturn through houses - complete 12 houses
		"saturn_house_1": "Period of increased self-discipline and responsibility. You'll take yourself and life goals more seriously, may experience some limits or challenges, but these are to build a stronger foundation. Important period for cultivating patience, taking responsibility, shaping mature personality.",
		"saturn_house_2": "**Financial responsibility** and **material management** are highlighted. You need **serious financial planning**, **stable income sources**; may face **financial limits** or need **spending restraint**. Good time to **build long-term financial security**, **learn value management**. Suitable for **budgeting**, **paying debt**, **investing in value assets**. Build **material foundation** through **discipline and patience**.",
		"saturn_house_3": "**Communication responsibility** and **learning discipline** are emphasized. You need **serious learning**, **improve communication**; may feel **expression limited** or need **more careful speech**. Good time to **build knowledge foundation**, **improve sibling relationships**. Suitable for **systematic learning**, **completing education**, **handling neighbor relations**. Improve **communication** through **patience and persistence**.",
		"saturn_house_4": "**Family responsibility** and **emotional maturity** are themes. You need **take family responsibility**, **handle family matters**; may face **family pressure** or need **build family structure**. Good time to **improve family relationships**, **handle property**, **build emotional security**. Suitable for **caring for family**, **handling family inheritance**, **building family traditions**. Build **stable family foundation** through **responsibility and commitment**.",
		"saturn_house_5": "**Creative responsibility** and **emotional discipline** are emphasized. You need **serious creation**, **take responsibility in relationships**; may feel **creativity limited** or need **more mature emotional expression**. Good time to **build creative habits**, **handle child relationships**, **learn emotional boundaries**. Suitable for **completing creative projects**, **educating children**, **building healthy entertainment habits**.",
		"saturn_house_6": "**Health responsibility** and **work discipline** are priorities. You need **serious health care**, **improve work habits**; may face **health challenges** or need **more regular life**. Good time to **build health habits**, **improve professional skills**, **handle work responsibility**. Suitable for **health checkups**, **health planning**, **improving work processes**. Build **healthy work-life balance** through **discipline and persistence**.",
		"saturn_house_7": "Period of relationship tests and maturity. Intimate relationships face reality tests; need serious commitment and responsibility. May bring relationship consolidation or clearing unhealthy connections. Suitable for building long-term stable partnerships, learning responsibility and boundaries in relationships.",
		"saturn_house_8": "**Deep responsibility** and **resource management** are emphasized. You need **serious shared resource handling**, **face deep fears**; may face **financial pressure** or need **handle debt**. Good time to **build resource management**, **handle inheritance**, **learn emotional depth**. Suitable for **paying debt**, **handling insurance**, **deep psychological work**.",
		"saturn_house_9": "**Belief responsibility** and **philosophical discipline** are themes. You need **serious faith**, **build philosophical system**; may feel **belief limited** or need **deeper learning**. Good time to **complete higher education**, **build belief system**, **handle legal matters**. Suitable for **deep learning**, **handling legal documents**, **building spiritual traditions**.",
		"saturn_house_10": "Critical period for career building. You'll face major career responsibilities and challenges, may need more effort. Important stage for building career reputation, achieving long-term career goals. Success requires patience, discipline, and persistent effort.",
		"saturn_house_11": "**Friendship responsibility** and **community discipline** are emphasized. You need **serious friendships**, **take responsibility in groups**; may feel **socially limited** or need **more mature community participation**. Good time to **build long-term friendships**, **establish group status**, **achieve collective goals**. Suitable for **meaningful group participation**, **building social networks**, **realizing long-term visions**.",
		"saturn_house_12": "**Inner responsibility** and **spiritual discipline** are themes. You need **face inner fears**, **handle karma**; may feel **need solitude** or need **deeper practice**. Good time for **inner healing**, **releasing old patterns**, **building spiritual discipline**. Suitable for **meditation**, **handling hidden issues**, **building inner structure**.",
		
		// Uranus through houses - complete 12 houses
		"uranus_house_1": "**Personal breakthrough** and **sudden change** are highlighted. You'll **suddenly change image**, **break personal limits**; **personal freedom** and **independence** increase. Good time for **personal breakthrough**, **image change**, **pursuing freedom**. Suitable for **personal change**, **breaking limits**, **pursuing independence**. Mind **change stability**; avoid **over-sudden**.",
		"uranus_house_2": "**Financial breakthrough** and **value change** are emphasized. You'll **suddenly change finances**, **break value concepts**; **financial freedom** and **value independence** increase. Good time for **financial breakthrough**, **value change**, **pursuing financial freedom**. Suitable for **financial innovation**, **value change**, **financial independence**. Mind **financial safety**; avoid **over-risk**.",
		"uranus_house_3": "**Communication breakthrough** and **thinking change** are themes. You'll **suddenly change thinking**, **break communication ways**; **thinking freedom** and **communication independence** increase. Good time for **thinking breakthrough**, **communication change**, **pursuing thinking freedom**. Suitable for **thinking innovation**, **communication change**, **learning new tech**.",
		"uranus_house_4": "**Family breakthrough** and **emotional change** are emphasized. You'll **suddenly change family**, **break family patterns**; **family freedom** and **emotional independence** increase. Good time for **family breakthrough**, **family change**, **pursuing family freedom**. Suitable for **family change**, **handling property**, **family independence**. Mind **family stability**; avoid **over-sudden**.",
		"uranus_house_5": "**Creative breakthrough** and **entertainment change** are highlighted. You'll **suddenly change creation**, **break entertainment ways**; **creative freedom** and **entertainment independence** increase. Good time for **creative breakthrough**, **entertainment change**, **pursuing creative freedom**. Suitable for **creative innovation**, **entertainment change**, **creative independence**.",
		"uranus_house_6": "**Work breakthrough** and **health change** are priorities. You'll **suddenly change work**, **break work ways**; **work freedom** and **health independence** increase. Good time for **work breakthrough**, **work change**, **pursuing work freedom**. Suitable for **work innovation**, **health change**, **work independence**. Mind **work stability**; avoid **over-sudden**.",
		"uranus_house_7": "**Relationship breakthrough** and **cooperation change** are emphasized. You'll **suddenly change relationships**, **break cooperation ways**; **relationship freedom** and **cooperation independence** increase. Good time for **relationship breakthrough**, **relationship change**, **pursuing relationship freedom**. Suitable for **relationship innovation**, **cooperation change**, **relationship independence**. Mind **relationship stability**; avoid **over-sudden**.",
		"uranus_house_8": "**Resource breakthrough** and **deep change** are themes. You'll **suddenly change resources**, **break deep relationships**; **resource freedom** and **deep independence** increase. Good time for **resource breakthrough**, **resource change**, **pursuing resource freedom**. Suitable for **resource innovation**, **investment change**, **resource independence**. Mind **resource safety**; avoid **over-risk**.",
		"uranus_house_9": "**Spiritual breakthrough** and **ideal change** are emphasized. You'll **suddenly change ideals**, **break spiritual ways**; **spiritual freedom** and **ideal independence** increase. Good time for **spiritual breakthrough**, **ideal change**, **pursuing spiritual freedom**. Suitable for **spiritual innovation**, **ideal change**, **spiritual independence**.",
		"uranus_house_10": "**Career breakthrough** and **career change** are highlighted. You'll **suddenly change career**, **break career ways**; **career freedom** and **career independence** increase. Good time for **career breakthrough**, **career change**, **pursuing career freedom**. Suitable for **career innovation**, **career change**, **career independence**. Mind **career stability**; avoid **over-sudden**.",
		"uranus_house_11": "**Community breakthrough** and **ideal change** are emphasized. You'll **suddenly change community**, **break ideal ways**; **community freedom** and **ideal independence** increase. Good time for **community breakthrough**, **community change**, **pursuing community freedom**. Suitable for **community innovation**, **ideal change**, **community independence**.",
		"uranus_house_12": "**Inner breakthrough** and **spiritual change** are themes. You'll **suddenly change inner**, **break spiritual ways**; **inner freedom** and **spiritual independence** increase. Good time for **inner breakthrough**, **inner change**, **pursuing inner freedom**. Suitable for **inner innovation**, **spiritual change**, **inner independence**.",
		
		// Neptune through houses - complete 12 houses
		"neptune_house_1": "**Personal ideal** and **spiritual expression** are highlighted. You'll **show spiritual charm**, **pursue ideal self**; **personal inspiration** and **spiritual expression** increase. Good time for **personal ideal**, **spiritual expression**, **pursuing ideal**. Suitable for **spiritual exploration**, **ideal expression**, **personal inspiration**. Mind **staying realistic**; avoid **over-idealization**.",
		"neptune_house_2": "**Financial ideal** and **value inspiration** are emphasized. You'll **idealize finances**, **pursue value inspiration**; **financial ideal** and **value inspiration** increase. Good time for **financial ideal**, **value inspiration**, **pursuing ideal value**. Suitable for **ideal investment**, **value inspiration**, **financial ideal**. Mind **financial reality**; avoid **over-idealization**.",
		"neptune_house_3": "**Communication ideal** and **thinking inspiration** are themes. You'll **idealize communication**, **pursue thinking inspiration**; **communication ideal** and **thinking inspiration** increase. Good time for **communication ideal**, **thinking inspiration**, **pursuing ideal communication**. Suitable for **ideal communication**, **thinking inspiration**, **communication ideal**. Mind **communication reality**; avoid **over-idealization**.",
		"neptune_house_4": "**Family ideal** and **emotional inspiration** are emphasized. You'll **idealize family**, **pursue emotional inspiration**; **family ideal** and **emotional inspiration** increase. Good time for **family ideal**, **emotional inspiration**, **pursuing ideal family**. Suitable for **ideal family**, **emotional inspiration**, **family ideal**. Mind **family reality**; avoid **over-idealization**.",
		"neptune_house_5": "**Creative ideal** and **entertainment inspiration** are highlighted. You'll **idealize creation**, **pursue entertainment inspiration**; **creative ideal** and **entertainment inspiration** increase. Good time for **creative ideal**, **entertainment inspiration**, **pursuing ideal creation**. Suitable for **ideal creation**, **entertainment inspiration**, **creative ideal**.",
		"neptune_house_6": "**Work ideal** and **health inspiration** are priorities. You'll **idealize work**, **pursue health inspiration**; **work ideal** and **health inspiration** increase. Good time for **work ideal**, **health inspiration**, **pursuing ideal work**. Suitable for **ideal work**, **health inspiration**, **work ideal**. Mind **work reality**; avoid **over-idealization**.",
		"neptune_house_7": "**Relationship ideal** and **cooperation inspiration** are emphasized. You'll **idealize relationships**, **pursue cooperation inspiration**; **relationship ideal** and **cooperation inspiration** increase. Good time for **relationship ideal**, **cooperation inspiration**, **pursuing ideal relationship**. Suitable for **ideal relationship**, **cooperation inspiration**, **relationship ideal**. Mind **relationship reality**; avoid **over-idealization**.",
		"neptune_house_8": "**Resource ideal** and **deep inspiration** are themes. You'll **idealize resources**, **pursue deep inspiration**; **resource ideal** and **deep inspiration** increase. Good time for **resource ideal**, **deep inspiration**, **pursuing ideal resource**. Suitable for **ideal resource**, **deep inspiration**, **resource ideal**. Mind **resource reality**; avoid **over-idealization**.",
		"neptune_house_9": "**Spiritual ideal** and **philosophical inspiration** are emphasized. You'll **idealize spirit**, **pursue philosophical inspiration**; **spiritual ideal** and **philosophical inspiration** increase. Good time for **spiritual ideal**, **philosophical inspiration**, **pursuing ideal spirit**. Suitable for **ideal spirit**, **philosophical inspiration**, **spiritual ideal**.",
		"neptune_house_10": "**Career ideal** and **career inspiration** are highlighted. You'll **idealize career**, **pursue career inspiration**; **career ideal** and **career inspiration** increase. Good time for **career ideal**, **career inspiration**, **pursuing ideal career**. Suitable for **ideal career**, **career inspiration**, **career ideal**. Mind **career reality**; avoid **over-idealization**.",
		"neptune_house_11": "**Community ideal** and **collective inspiration** are emphasized. You'll **idealize community**, **pursue collective inspiration**; **community ideal** and **collective inspiration** increase. Good time for **community ideal**, **collective inspiration**, **pursuing ideal community**. Suitable for **ideal community**, **collective inspiration**, **community ideal**.",
		"neptune_house_12": "**Inner ideal** and **spiritual inspiration** are themes. You'll **idealize inner**, **pursue spiritual inspiration**; **inner ideal** and **spiritual inspiration** increase. Good time for **inner ideal**, **spiritual inspiration**, **pursuing ideal inner**. Suitable for **ideal inner**, **spiritual inspiration**, **inner ideal**.",
		
		// Pluto through houses - complete 12 houses
		"pluto_house_1": "**Personal transformation** and **deep change** are highlighted. You'll **experience deep transformation**, **change personal essence**; **personal power** and **transformation ability** increase. Good time for **personal transformation**, **deep change**, **pursuing personal power**. Suitable for **personal transformation**, **deep change**, **personal power**. Mind **transformation stability**; avoid **over-intense**.",
		"pluto_house_2": "**Financial transformation** and **value change** are emphasized. You'll **experience financial transformation**, **change value concepts**; **financial power** and **value transformation** increase. Good time for **financial transformation**, **value change**, **pursuing financial power**. Suitable for **financial transformation**, **value change**, **financial power**. Mind **financial safety**; avoid **over-intense**.",
		"pluto_house_3": "**Communication transformation** and **thinking change** are themes. You'll **experience communication transformation**, **change thinking ways**; **communication power** and **thinking transformation** increase. Good time for **communication transformation**, **thinking change**, **pursuing communication power**. Suitable for **communication transformation**, **thinking change**, **communication power**.",
		"pluto_house_4": "**Family transformation** and **emotional change** are emphasized. You'll **experience family transformation**, **change family patterns**; **family power** and **emotional transformation** increase. Good time for **family transformation**, **emotional change**, **pursuing family power**. Suitable for **family transformation**, **emotional change**, **family power**. Mind **family stability**; avoid **over-intense**.",
		"pluto_house_5": "**Creative transformation** and **entertainment change** are highlighted. You'll **experience creative transformation**, **change entertainment ways**; **creative power** and **entertainment transformation** increase. Good time for **creative transformation**, **entertainment change**, **pursuing creative power**. Suitable for **creative transformation**, **entertainment change**, **creative power**.",
		"pluto_house_6": "**Work transformation** and **health change** are priorities. You'll **experience work transformation**, **change work ways**; **work power** and **health transformation** increase. Good time for **work transformation**, **health change**, **pursuing work power**. Suitable for **work transformation**, **health change**, **work power**. Mind **work stability**; avoid **over-intense**.",
		"pluto_house_7": "**Relationship transformation** and **cooperation change** are emphasized. You'll **experience relationship transformation**, **change cooperation ways**; **relationship power** and **cooperation transformation** increase. Good time for **relationship transformation**, **cooperation change**, **pursuing relationship power**. Suitable for **relationship transformation**, **cooperation change**, **relationship power**. Mind **relationship stability**; avoid **over-intense**.",
		"pluto_house_8": "**Resource transformation** and **deep change** are themes. You'll **experience resource transformation**, **change deep relationships**; **resource power** and **deep transformation** increase. Good time for **resource transformation**, **deep change**, **pursuing resource power**. Suitable for **resource transformation**, **investment change**, **resource power**. Mind **resource safety**; avoid **over-intense**.",
		"pluto_house_9": "**Spiritual transformation** and **ideal change** are emphasized. You'll **experience spiritual transformation**, **change ideal ways**; **spiritual power** and **ideal transformation** increase. Good time for **spiritual transformation**, **ideal change**, **pursuing spiritual power**. Suitable for **spiritual transformation**, **ideal change**, **spiritual power**.",
		"pluto_house_10": "**Career transformation** and **career change** are highlighted. You'll **experience career transformation**, **change career ways**; **career power** and **career transformation** increase. Good time for **career transformation**, **career change**, **pursuing career power**. Suitable for **career transformation**, **career change**, **career power**. Mind **career stability**; avoid **over-intense**.",
		"pluto_house_11": "**Community transformation** and **ideal change** are emphasized. You'll **experience community transformation**, **change ideal ways**; **community power** and **ideal transformation** increase. Good time for **community transformation**, **ideal change**, **pursuing community power**. Suitable for **community transformation**, **ideal change**, **community power**.",
		"pluto_house_12": "**Inner transformation** and **spiritual change** are themes. You'll **experience inner transformation**, **change spiritual ways**; **inner power** and **spiritual transformation** increase. Good time for **inner transformation**, **spiritual change**, **pursuing inner power**. Suitable for **inner transformation**, **spiritual change**, **inner power**.",
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
		// Sun through houses - complete 12 houses
		"sun_house_1": "**Личное обаяние** и **жизненная энергия** полностью усилены. Вы почувствуете себя **физически энергичным** с отличным ментальным состоянием, **выносливость и сила** на пике. Эта мощная жизненная сила позволяет вам **преуспевать в профессиональной среде**, **лидерские качества и личный магнетизм** привлекают внимание. Хорошее время для проявления инициативы на **работе**, демонстрируя вашу **профессиональную компетентность и здоровую жизнеспособность**.",
		"sun_house_2": "**Финансовые вопросы** и **материальная безопасность** становятся фокусом. Вы уделите больше внимания **источникам дохода** и **статусу активов**, с возможными возможностями **увеличить заработок** или пересмотреть **финансовое планирование**. Более четкое понимание вашей ценности повышает **способность зарабатывать**. Подходит для **финансового планирования**, **инвестиций**, **управления активами**. Обратите внимание на **рациональное управление деньгами**, избегайте **чрезмерных трат**.",
		"sun_house_3": "**Общение** приносит **карьерные возможности** и **расширение межличностных связей**. Вы обнаружите, что **сотрудничество на рабочем месте** становится частым с большим количеством **деловых встреч и презентаций**. Между тем **социальные сети** расширяются, так как **заводить новых друзей и поддерживать старые отношения** идет гладко. Хорошее время для **построения профессиональных связей** и **продвижения проектов через коммуникативные навыки**. Подходит для **обучения**, **письма**, **коротких поездок**. Обратите внимание на **точность информации**, избегайте **распространения ошибок**.",
		"sun_house_4": "**Семейные отношения** и **эмоциональные связи** выделены. Вы уделите больше внимания **домашней жизни** и **семейным делам**, **семейная гармония** и **эмоциональная безопасность** улучшаются. Хорошее время для **улучшения семейных отношений**, **обработки недвижимости**, **создания семейной традиции**. Подходит для **семейных мероприятий**, **обработки недвижимости**, **семейного времени**. Обратите внимание на **семейную гармонию**, избегайте **семейных конфликтов**.",
		"sun_house_5": "**Творческое выражение** и **развлечения** подчеркнуты. Вы будете **активно творить**, **наслаждаться развлечениями**, **творческие способности** и **радость жизни** усиливаются. Хорошее время для **творческих проектов**, **развлечений**, **образования детей**. Подходит для **художественного творчества**, **спорта**, **развлечений**. Обратите внимание на **творческую свободу**, избегайте **чрезмерного контроля**.",
		"sun_house_6": "**Работа** и **здоровье** становятся приоритетом. Вы будете **активно работать**, **улучшать здоровье**, **рабочая эффективность** и **здоровье** улучшаются. Хорошее время для **улучшения рабочего процесса**, **управления здоровьем**, **обработки рабочих задач**. Подходит для **рабочих проектов**, **фитнеса**, **улучшения здоровья**. Обратите внимание на **баланс работы и жизни**, избегайте **переутомления**.",
		"sun_house_7": "**Партнерские отношения** и **сотрудничество** становятся фокусом. Вы будете **уделять больше внимания отношениям**, **улучшать партнерство**, **качество отношений** и **способность к сотрудничеству** улучшаются. Хорошее время для **улучшения отношений**, **заключения контрактов**, **обработки партнерства**. Подходит для **отношений**, **бизнес-переговоров**, **создания партнерства**. Обратите внимание на **баланс в отношениях**, избегайте **одностороннего доминирования**.",
		"sun_house_8": "**Глубокие отношения** и **общие ресурсы** подчеркнуты. Вы будете **углубляться в отношения**, **обрабатывать общие ресурсы**, **глубина отношений** и **управление ресурсами** улучшаются. Хорошее время для **углубления отношений**, **обработки инвестиций**, **управления общими ресурсами**. Подходит для **финансовых исследований**, **инвестиций**, **обработки наследства**. Обратите внимание на **безопасность ресурсов**, избегайте **чрезмерного риска**.",
		"sun_house_9": "**Философское мышление** и **духовное развитие** достигают пика. Вы будете **размышлять о смысле жизни**, **изучать философию**, **высшее образование** и **путешествия** идут гладко. Хорошее время для **обучения**, **преподавания**, **путешествий**, **юридических дел**. Подходит для **высшего образования**, **публикаций**, **юридических консультаций**, **межкультурного обучения**.",
		"sun_house_10": "**Развитие карьеры** и **профессиональные достижения** становятся центром внимания. Ваш **профессиональный путь** и **социальный статус** подчеркнуты, возможны **возможности повышения** или больше **рабочих обязанностей**. Ваши **профессиональные способности** и **результаты работы** получают общественное признание. Критический период для построения **профессиональной репутации** и реализации **карьерных амбиций**. Подходит для **профессионального развития**, **публичных выступлений**, **профессиональной сертификации**. Обратите внимание на **поддержание профессионального имиджа**.",
		"sun_house_11": "**Социальные связи** и **коллективные цели** подчеркнуты. Вы будете **активно участвовать в группах**, **обсуждать коллективные видения**, **социальные сети** и **коллективные цели** выигрывают. Хорошее время для **участия в групповых мероприятиях**, **создания социальных сетей**, **распространения идей**. Подходит для **социальных сетей**, **групповых обсуждений**, **реализации коллективных видений**.",
		"sun_house_12": "**Внутреннее размышление** и **духовное исследование** становятся темой. Вы будете **глубоко размышлять**, **исследовать подсознание**, **интуиция** и **вдохновение** становятся активными. Хорошее время для **внутреннего обучения**, **духовного исследования**, **обработки скрытых проблем**. Подходит для **медитации**, **письма**, **психологического исцеления**. Обратите внимание на **ясность информации**, избегайте **путаницы в мыслях**.",
		
		// Moon through houses - complete 12 houses
		"moon_house_1": "**Эмоциональное выражение** и **личная чувствительность** становятся фокусом. Вы будете **более эмоционально выражаться**, **проявлять личную чувствительность**, **эмоциональная осведомленность** и **интуиция** усиливаются. Хорошее время для **эмоционального выражения**, **личного роста**, **улучшения эмоционального здоровья**. Подходит для **эмоционального исцеления**, **самоухода**, **эмоционального развития**. Обратите внимание на **эмоциональный баланс**, избегайте **чрезмерной чувствительности**.",
		"moon_house_2": "**Эмоциональная безопасность** и **материальные ценности** подчеркнуты. Вы будете **искать эмоциональную безопасность через материальные ценности**, **управлять финансами с эмоциональной точки зрения**, **финансовая безопасность** и **эмоциональная стабильность** связаны. Хорошее время для **финансового планирования**, **управления активами**, **построения материальной безопасности**. Подходит для **финансового управления**, **инвестиций**, **построения финансовой стабильности**. Обратите внимание на **рациональное управление деньгами**, избегайте **эмоциональных трат**.",
		"moon_house_3": "**Эмоциональное общение** и **интуитивное обучение** достигают пика. Вы будете **общаться с эмоциями**, **учиться через интуицию**, **эмоциональное общение** и **интуитивное обучение** улучшаются. Хорошее время для **эмоционального общения**, **интуитивного обучения**, **коротких поездок**. Подходит для **общения с семьей**, **обучения**, **коротких поездок**. Обратите внимание на **эмоциональную ясность**, избегайте **эмоциональной путаницы**.",
		"moon_house_4": "**Семейные эмоции** и **эмоциональная безопасность** становятся темой. Вы будете **больше времени проводить с семьей**, **обрабатывать семейные эмоции**, **семейная гармония** и **эмоциональная безопасность** улучшаются. Хорошее время для **улучшения семейных отношений**, **обработки семейных дел**, **построения эмоциональной безопасности**. Подходит для **семейных мероприятий**, **обработки недвижимости**, **семейного времени**. Обратите внимание на **семейную гармонию**, избегайте **семейных конфликтов**.",
		"moon_house_5": "**Эмоциональное творчество** и **радость жизни** подчеркнуты. Вы будете **творить с эмоциями**, **наслаждаться жизнью**, **творческие способности** и **радость жизни** усиливаются. Хорошее время для **творческих проектов**, **развлечений**, **образования детей**. Подходит для **художественного творчества**, **развлечений**, **образования детей**. Обратите внимание на **эмоциональную свободу**, избегайте **эмоционального контроля**.",
		"moon_house_6": "**Эмоциональное здоровье** и **рабочие эмоции** становятся приоритетом. Вы будете **обрабатывать эмоциональное здоровье**, **управлять рабочими эмоциями**, **здоровье** и **рабочая эффективность** улучшаются. Хорошее время для **улучшения здоровья**, **управления рабочими эмоциями**, **обработки рабочих задач**. Подходит для **здорового образа жизни**, **рабочих проектов**, **улучшения здоровья**. Обратите внимание на **эмоциональное здоровье**, избегайте **эмоционального стресса**.",
		"moon_house_7": "**Эмоциональные отношения** и **партнерские эмоции** становятся фокусом. Вы будете **уделять больше внимания эмоциональным отношениям**, **улучшать партнерские эмоции**, **качество отношений** и **эмоциональная связь** улучшаются. Хорошее время для **улучшения эмоциональных отношений**, **обработки партнерства**, **построения эмоциональной связи**. Подходит для **отношений**, **эмоционального общения**, **построения отношений**. Обратите внимание на **эмоциональный баланс**, избегайте **эмоциональной зависимости**.",
		"moon_house_8": "**Глубокие эмоции** и **эмоциональная трансформация** подчеркнуты. Вы будете **углубляться в эмоции**, **обрабатывать эмоциональную трансформацию**, **глубина эмоций** и **эмоциональная трансформация** улучшаются. Хорошее время для **углубления эмоций**, **обработки эмоциональной трансформации**, **управления общими ресурсами**. Подходит для **эмоционального исцеления**, **инвестиций**, **обработки наследства**. Обратите внимание на **эмоциональную безопасность**, избегайте **эмоциональной манипуляции**.",
		"moon_house_9": "**Эмоциональная философия** и **духовные эмоции** достигают пика. Вы будете **размышлять о жизни с эмоциями**, **изучать духовность**, **эмоциональная философия** и **духовные эмоции** улучшаются. Хорошее время для **эмоционального обучения**, **духовного исследования**, **путешествий**. Подходит для **обучения**, **путешествий**, **духовного исследования**. Обратите внимание на **эмоциональную ясность**, избегайте **эмоциональной путаницы**.",
		"moon_house_10": "**Эмоциональная карьера** и **публичные эмоции** становятся фокусом. Вы будете **выражать эмоции в карьере**, **управлять публичными эмоциями**, **карьера** и **публичный имидж** выигрывают. Хорошее время для **профессионального развития**, **публичных выступлений**, **построения профессиональной репутации**. Подходит для **профессионального развития**, **публичных выступлений**, **профессиональной сертификации**. Обратите внимание на **эмоциональный баланс**, избегайте **эмоциональной нестабильности**.",
		"moon_house_11": "**Эмоциональные дружеские отношения** и **коллективные эмоции** подчеркнуты. Вы будете **строить эмоциональные дружеские отношения**, **участвовать в коллективных эмоциях**, **социальные сети** и **коллективные цели** выигрывают. Хорошее время для **участия в групповых мероприятиях**, **создания социальных сетей**, **распространения идей**. Подходит для **социальных сетей**, **групповых обсуждений**, **реализации коллективных видений**.",
		"moon_house_12": "Эмоции входят в подсознательный уровень, вы можете почувствовать эмоциональную неопределенность и чувствительность. Нужно побыть в одиночестве, чтобы обработать внутренние эмоции, сны и интуиция становятся активными. Это время подходит для **медитации**, **одиночества**, **эмоционального исцеления**. Вы можете почувствовать эмоциональную уязвимость, нужно дать себе пространство и время для восстановления.",
		
		// Mercury through houses - complete 12 houses
		"mercury_house_1": "**Мыслительное выражение** и **личное общение** становятся фокусом. Вы будете **чаще выражать мнения**, **проявлять личную мудрость**, **коммуникативные способности** и **способности к обучению** улучшаются. Это хорошее время для **проявления личного обаяния**, **создания личного бренда**, **улучшения способности к выражению**. Подходит для **публичных выступлений**, **письма**, **обучения новым навыкам**. Обратите внимание на **точность общения**, избегайте **слишком поспешного или рассеянного**.",
		"mercury_house_2": "**Финансовое общение** и **мышление о ценностях** подчеркнуты. Вы будете **думать о финансовых вопросах**, **обсуждать денежные темы**, **финансовые концепции** и **суждения о ценностях** становятся важными. Это хорошее время для **финансового планирования**, **обучения финансовым знаниям**, **обсуждения условий контракта**. Подходит для **финансовых переговоров**, **оценки активов**, **обучения инвестированию**. Обратите внимание на **точность финансовой информации**, избегайте **импульсивных решений**.",
		"mercury_house_3": "**Общение** и **обучающее выражение** достигают пика. Вы будете **часто общаться**, **совершать короткие поездки**, **обучение** и **письмо** идут гладко. Это отличное время для **улучшения коммуникативных навыков**, **создания сетей**, **обучения новым знаниям**. Подходит для **встреч**, **письма**, **преподавания**, **коротких поездок**. Обратите внимание на **точность информации**, избегайте **распространения ошибочной информации**.",
		"mercury_house_4": "**Семейное общение** и **эмоциональное выражение** становятся темой. Вы будете **больше общаться с семьей**, **обсуждать семейные дела**, **семейные отношения** и **эмоциональное выражение** становятся важными. Это хорошее время для **улучшения семейного общения**, **обработки семейных документов**, **изучения семейной истории**. Подходит для **семейных встреч**, **обработки документов на недвижимость**, **записи семейных историй**.",
		"mercury_house_5": "**Творческое выражение** и **развлекательное общение** подчеркнуты. Вы будете **выражаться через творчество**, **общаться в развлечениях**, **творческие способности** и **способности к выражению** объединяются. Это хорошее время для **письма**, **преподавания**, **творческих проектов**. Подходит для **художественного творчества**, **образования детей**, **развлечений**. Обратите внимание на **интересность выражения**, избегайте **слишком серьезного**.",
		"mercury_house_6": "**Рабочее общение** и **обучение здоровью** становятся приоритетом. Вы будете **часто общаться на работе**, **изучать знания о здоровье**, **рабочая эффективность** и **управление здоровьем** выигрывают. Это хорошее время для **улучшения рабочего процесса**, **обучения знаниям о здоровье**, **обработки рабочих документов**. Подходит для **рабочего обучения**, **консультаций по здоровью**, **оптимизации рабочих методов**.",
		"mercury_house_7": "**Общение в отношениях** и **кооперативное общение** становятся фокусом. Вы будете **часто общаться с партнерами**, **обсуждать кооперативные отношения**, **коммуникативные навыки** в **отношениях** становятся важными. Это хорошее время для **улучшения общения в отношениях**, **заключения контрактов**, **обработки кооперативных отношений**. Подходит для **консультаций по отношениям**, **бизнес-переговоров**, **создания кооперации**. Обратите внимание на **баланс общения**, избегайте **одностороннего доминирования**.",
		"mercury_house_8": "**Глубокое общение** и **исследование ресурсов** подчеркнуты. Вы будете **глубоко обсуждать**, **исследовать вопросы ресурсов**, **психологическая проницательность** и **способности к анализу ресурсов** улучшаются. Это хорошее время для **глубокого исследования**, **обработки общих ресурсов**, **изучения психологии**. Подходит для **финансовых исследований**, **психологических консультаций**, **обработки наследства**. Обратите внимание на **конфиденциальность информации**, избегайте **разглашения конфиденциальной информации**.",
		"mercury_house_9": "**Философское мышление** и **распространение мудрости** достигают пика. Вы будете **размышлять о смысле жизни**, **изучать философию**, **высшее образование** и **путешествия** идут гладко. Это отличное время для **обучения**, **преподавания**, **путешествий**, **юридических дел**. Подходит для **высшего образования**, **публикаций**, **юридических консультаций**, **межкультурного обучения**.",
		"mercury_house_10": "**Профессиональное общение** и **публичное выражение** становятся фокусом. Вы будете **часто общаться в профессии**, **проявлять профессиональные способности**, **профессиональная репутация** и **публичный имидж** выигрывают. Это хорошее время для **профессиональных выступлений**, **профессионального письма**, **создания профессиональных сетей**. Подходит для **профессионального развития**, **публичных выступлений**, **профессиональной сертификации**. Обратите внимание на **поддержание профессионального имиджа**.",
		"mercury_house_11": "**Социальное общение** и **идеальное общение** подчеркнуты. Вы будете **активно общаться в группах**, **обсуждать идеальные видения**, **социальные сети** и **коллективные цели** выигрывают. Это хорошее время для **участия в групповых мероприятиях**, **создания социальных сетей**, **распространения идей**. Подходит для **социальных сетей**, **групповых обсуждений**, **реализации коллективных видений**.",
		"mercury_house_12": "**Внутреннее мышление** и **подсознательное общение** становятся темой. Вы будете **глубоко размышлять**, **исследовать подсознание**, **интуиция** и **вдохновение** становятся активными. Это хорошее время для **внутреннего обучения**, **духовного исследования**, **обработки скрытых проблем**. Подходит для **медитации**, **письма**, **психологического исцеления**. Обратите внимание на **ясность информации**, избегайте **путаницы в мыслях**.",
		
		// Venus through houses - complete 12 houses
		"venus_house_1": "**Личное обаяние** и **эстетическое улучшение** становятся фокусом. Вы будете **больше внимания уделять внешности**, **проявлять личное обаяние**, **привлекательность** и **популярность** улучшаются. Это хорошее время для **улучшения образа**, **повышения обаяния**, **наслаждения жизнью**. Подходит для **красоты**, **покупок**, **общения**. Обратите внимание на **сохранение подлинности**, избегайте **чрезмерного внимания к внешности**.",
		"venus_house_2": "**Финансовое наслаждение** и **материальная красота** подчеркнуты. Вы будете **наслаждаться материальной жизнью**, **инвестировать в красивые вещи**, **доход** и **расходы** выигрывают. Это хорошее время для **увеличения дохода**, **наслаждения материальным**, **инвестирования в ценности**. Подходит для **финансового управления**, **покупок**, **инвестирования в искусство**. Обратите внимание на **рациональные расходы**, избегайте **чрезмерных трат**.",
		"venus_house_3": "**Обаяние в общении** и **социальное удовольствие** становятся темой. Вы будете **проявлять обаяние в общении**, **наслаждаться общением**, **межличностные отношения** и **обучение** выигрывают. Это хорошее время для **улучшения общения**, **создания сетей**, **наслаждения обучением**. Подходит для **общения**, **обучения**, **коротких поездок**.",
		"venus_house_4": "**Семейная гармония** и **эмоциональная красота** подчеркнуты. Вы будете **украшать дом**, **наслаждаться семейным временем**, **семейные отношения** и **эмоциональная безопасность** выигрывают. Это хорошее время для **улучшения дома**, **семейных собраний**, **наслаждения семейной жизнью**. Подходит для **домашнего декора**, **семейных мероприятий**, **обработки недвижимости**.",
		"venus_house_5": "Романтика и творчество достигают пика. Этот период полон **обаяния и привлекательности**, романтические отношения могут начаться или усилиться. Вы будете **больше наслаждаться радостями жизни**, художественное творчество и развлечения приносят удовольствие. Взаимодействие с детьми также полно радости. Это прекрасное время для **погони за романтикой**, **выражения любви**, **наслаждения творческой деятельностью**.",
		"venus_house_6": "**Красота работы** и **наслаждение здоровьем** становятся приоритетом. Вы будете **украшать рабочую среду**, **наслаждаться здоровой жизнью**, **рабочие отношения** и **физическое здоровье** выигрывают. Это хорошее время для **улучшения рабочей среды**, **наслаждения работой**, **здоровой красоты**. Подходит для **рабочего общения**, **управления здоровьем**, **наслаждения повседневным**.",
		"venus_house_7": "Гармония отношений становится фокусом. Вы будете **больше ценить партнерские отношения**, **стремиться к балансу и красивым связям**. Это хорошее время для **улучшения интимных отношений**, **создания партнерских отношений**. Ваши **дипломатические навыки** и **обаяние** помогают вам преуспевать в межличностных отношениях. Подходит для **решения конфликтов в отношениях**, **заключения контрактов**.",
		"venus_house_8": "**Глубокая привлекательность** и **красота ресурсов** подчеркнуты. Вы будете **углубляться в отношения**, **наслаждаться общими ресурсами**, **интимные отношения** и **финансовое разделение** выигрывают. Это хорошее время для **углубления отношений**, **обработки общего имущества**, **наслаждения ресурсами**. Подходит для **углубления отношений**, **инвестиций**, **обработки наследства**. Обратите внимание на **глубину отношений**, избегайте **чрезмерной зависимости**.",
		"venus_house_9": "**Духовная красота** и **наслаждение идеалами** становятся темой. Вы будете **наслаждаться духовными поисками**, **изучать эстетику**, **вера** и **путешествия** выигрывают. Это хорошее время для **изучения искусства**, **духовных путешествий**, **наслаждения культурой**. Подходит для **изучения искусства**, **культурных путешествий**, **духовного исследования**.",
		"venus_house_10": "**Профессиональное обаяние** и **публичная красота** становятся фокусом. Вы будете **проявлять обаяние в профессии**, **наслаждаться профессиональными достижениями**, **профессиональные отношения** и **публичный имидж** выигрывают. Это хорошее время для **повышения профессионального образа**, **наслаждения профессиональным успехом**, **создания профессиональных отношений**. Подходит для **профессионального общения**, **публичного образа**, **профессионального развития**.",
		"venus_house_11": "Расцвет социальной харизмы - дружба и общественная деятельность приносят радость. Вы заведете новых друзей, будете участвовать в интересных групповых мероприятиях. Ваша социальная сеть расширяется, возможно, вы встретите романтический объект через круг друзей. Это время подходит для **участия в социальных мероприятиях**, **погони за коллективными целями**, **наслаждения красотой дружбы**.",
		"venus_house_12": "**Внутренняя красота** и **духовное наслаждение** подчеркнуты. Вы будете **наслаждаться одиночеством**, **исследовать внутреннюю красоту**, **духовность** и **искусство** выигрывают. Это хорошее время для **внутренней красоты**, **духовного искусства**, **наслаждения покоем**. Подходит для **медитации**, **художественного творчества**, **внутреннего исследования**.",
		
		// Mars through houses - complete 12 houses
		"mars_house_1": "Действенность и личная инициатива значительно усиливаются. Вы почувствуете себя **энергичным**, захотите начать новые проекты или вызовы. **Уверенность в себе** повышается, но также можете стать **импульсивным**. Это хорошее время для **достижения личных целей**, **проявления лидерства**, но обратите внимание на **контроль темперамента**, избегайте **конфликтов с другими**.",
		"mars_house_2": "**Финансовые действия** и **материальная конкуренция** подчеркнуты. Вы будете **активно стремиться к доходу**, **конкурировать за ресурсы**, **способность зарабатывать** и **импульс к расходам** усиливаются. Это хорошее время для **увеличения дохода**, **инвестиций**, **достижения финансовых целей**. Подходит для **финансовых действий**, **инвестиций**, **борьбы за повышение зарплаты**. Обратите внимание на **рациональные расходы**, избегайте **импульсивных покупок**.",
		"mars_house_3": "**Действия в общении** и **конкуренция в обучении** становятся темой. Вы будете **активно общаться**, **конкурировать в обучении**, **коммуникативные способности** и **мотивация к обучению** улучшаются. Это хорошее время для **активного общения**, **конкуренции в обучении**, **коротких поездок**. Подходит для **дебатов**, **обучения**, **коротких поездок**. Обратите внимание на **тон общения**, избегайте **слишком агрессивного**.",
		"mars_house_4": "**Семейные действия** и **эмоциональная конкуренция** подчеркнуты. Вы будете **активно обрабатывать семейные дела**, **бороться за семейное положение**, **семейные отношения** и **эмоциональное выражение** становятся интенсивными. Это хорошее время для **улучшения семьи**, **обработки семейных конфликтов**, **достижения семейных целей**. Подходит для **семейных проектов**, **обработки недвижимости**, **семейных мероприятий**. Обратите внимание на **семейную гармонию**, избегайте **семейных конфликтов**.",
		"mars_house_5": "**Творческие действия** и **развлекательная конкуренция** становятся фокусом. Вы будете **активно творить**, **конкурировать в развлечениях**, **творческие способности** и **действенность** объединяются. Это хорошее время для **активного творчества**, **спорта**, **конкуренции в развлечениях**. Подходит для **художественного творчества**, **спорта**, **развлекательной конкуренции**. Обратите внимание на **безопасность**, избегайте **чрезмерного риска**.",
		"mars_house_6": "**Рабочие действия** и **конкуренция в здоровье** становятся приоритетом. Вы будете **активно работать**, **конкурировать в здоровье**, **рабочая эффективность** и **управление здоровьем** выигрывают. Это хорошее время для **активной работы**, **спорта и фитнеса**, **конкуренции за рабочие цели**. Подходит для **рабочих спринтов**, **фитнеса**, **улучшения здоровья**. Обратите внимание на **баланс работы**, избегайте **переутомления**.",
		"mars_house_7": "**Действия в отношениях** и **конкуренция в сотрудничестве** подчеркнуты. Вы будете **активно обрабатывать отношения**, **конкурировать в сотрудничестве**, **динамика отношений** и **способность к сотрудничеству** улучшаются. Это хорошее время для **активного улучшения отношений**, **конкуренции в сотрудничестве**, **обработки конфликтов в отношениях**. Подходит для **действий в отношениях**, **конкуренции в сотрудничестве**, **обработки отношений**. Обратите внимание на **баланс в отношениях**, избегайте **конфликтов в отношениях**.",
		"mars_house_8": "**Глубокие действия** и **конкуренция за ресурсы** становятся темой. Вы будете **активно обрабатывать ресурсы**, **конкурировать за глубокие отношения**, **управление ресурсами** и **эмоциональная глубина** становятся интенсивными. Это хорошее время для **активной обработки ресурсов**, **конкуренции за инвестиции**, **обработки общего имущества**. Подходит для **действий с ресурсами**, **конкуренции за инвестиции**, **обработки долгов**. Обратите внимание на **безопасность ресурсов**, избегайте **чрезмерного риска**.",
		"mars_house_9": "**Духовные действия** и **конкуренция за идеалы** подчеркнуты. Вы будете **активно стремиться к идеалам**, **конкурировать за духовные цели**, **духовные поиски** и **реализация идеалов** выигрывают. Это хорошее время для **активного обучения**, **конкуренции в духовном**, **погони за идеалами**. Подходит для **конкуренции в обучении**, **духовных поисков**, **действий за идеалы**.",
		"mars_house_10": "Карьерные амбиции и профессиональная движущая сила достигают пика. Вы будете **всеми силами стремиться к карьерным целям**, могут быть важные **карьерные прорывы**. На работе проявляется **мощная исполнительная способность** и **конкурентоспособность**. Это хорошее время для **борьбы за продвижение**, **продвижения важных проектов**, но обратите внимание на **отношения с авторитетными фигурами**.",
		"mars_house_11": "Активно и смело стремитесь к идеалам и социальным целям. Вы будете **участвовать в групповых мероприятиях**, **бороться за общие цели**. В общении проявляется **действенность** и **лидерские способности**. Это время подходит для **продвижения социальных изменений**, **реализации коллективных целей**, но обратите внимание на **командную работу**, избегайте **самовольства**.",
		"mars_house_12": "**Внутренние действия** и **духовная конкуренция** становятся темой. Вы будете **активно обрабатывать внутренние проблемы**, **конкурировать за духовные цели**, **внутренняя динамика** и **духовные поиски** выигрывают. Это хорошее время для **активного исцеления**, **внутренних действий**, **духовной конкуренции**. Подходит для **внутренней работы**, **духовных действий**, **обработки скрытых проблем**. Обратите внимание на **внутренний баланс**, избегайте **внутренних конфликтов**.",
		
		// Jupiter through houses - complete 12 houses
		"jupiter_house_1": "**Личное расширение** и **оптимистичная уверенность** становятся фокусом. Вы будете **чувствовать уверенность**, **проявлять личное обаяние**, **личный рост** и **возможности** увеличиваются. Это хорошее время для **личного развития**, **проявления обаяния**, **захвата возможностей**. Подходит для **личных целей**, **проявления способностей**, **захвата возможностей**. Обратите внимание на **сохранение скромности**, избегайте **чрезмерной уверенности**.",
		"jupiter_house_2": "Период увеличения богатства и возможностей. Доход может повыситься или появятся новые возможности заработать. Ваше отношение к материальному становится более **оптимистичным и щедрым**. Это хорошее время для **инвестиций**, **расширения финансовых ресурсов**. Обратите внимание на то, чтобы не **чрезмерно тратить** или быть **слишком оптимистичным**, сохраняйте **рациональное финансовое планирование**.",
		"jupiter_house_3": "**Расширение общения** и **возможности обучения** подчеркнуты. Вы будете **часто общаться**, **возможности обучения** увеличиваются, **коммуникативные способности** и **способности к обучению** улучшаются. Это хорошее время для **обучения**, **общения**, **коротких поездок**. Подходит для **обучения**, **преподавания**, **создания сетей**.",
		"jupiter_house_4": "**Семейное расширение** и **эмоциональный оптимизм** становятся темой. Вы будете **улучшать семейные отношения**, **семейные возможности** увеличиваются, **семейные отношения** и **эмоциональная безопасность** выигрывают. Это хорошее время для **улучшения семьи**, **семейного расширения**, **эмоционального оптимизма**. Подходит для **семейных мероприятий**, **обработки недвижимости**, **семейного роста**.",
		"jupiter_house_5": "**Творческое расширение** и **возможности развлечений** подчеркнуты. Вы будете **увеличивать творческие возможности**, **повышать наслаждение развлечениями**, **творческие способности** и **способности к развлечениям** выигрывают. Это хорошее время для **творчества**, **развлечений**, **образования детей**. Подходит для **художественного творчества**, **развлечений**, **образования**.",
		"jupiter_house_6": "**Рабочее расширение** и **возможности здоровья** становятся приоритетом. Вы будете **увеличивать рабочие возможности**, **улучшать здоровье**, **рабочая эффективность** и **управление здоровьем** выигрывают. Это хорошее время для **рабочего развития**, **улучшения здоровья**, **рабочих возможностей**. Подходит для **рабочего расширения**, **управления здоровьем**, **рабочего обучения**.",
		"jupiter_house_7": "**Расширение отношений** и **возможности сотрудничества** становятся фокусом. Вы будете **увеличивать возможности отношений**, **расширять сотрудничество**, **качество отношений** и **способность к сотрудничеству** выигрывают. Это хорошее время для **развития отношений**, **расширения сотрудничества**, **возможностей отношений**. Подходит для **развития отношений**, **сотрудничества**, **обучения отношениям**.",
		"jupiter_house_8": "**Расширение ресурсов** и **глубокие возможности** подчеркнуты. Вы будете **увеличивать возможности ресурсов**, **расширять глубокие отношения**, **управление ресурсами** и **эмоциональная глубина** выигрывают. Это хорошее время для **расширения ресурсов**, **инвестиционных возможностей**, **глубоких отношений**. Подходит для **инвестиций**, **управления ресурсами**, **глубоких отношений**.",
		"jupiter_house_9": "Золотой период мудрого просвещения и духовного роста. Вы будете **интересоваться философией**, **религией**, **высшим образованием**. Могут быть возможности **дальних поездок** или **контакта с разными культурами**. **Видение расширяется**, понимание смысла жизни становится глубже. Это отличное время для **обучения**, **путешествий**, **погони за истиной**.",
		"jupiter_house_10": "**Профессиональное расширение** и **карьерные возможности** становятся фокусом. Вы будете **увеличивать профессиональные возможности**, **расширять карьеру**, **профессиональное развитие** и **карьерные достижения** выигрывают. Это хорошее время для **профессионального развития**, **расширения карьеры**, **профессиональных возможностей**. Подходит для **профессионального развития**, **расширения карьеры**, **профессионального обучения**.",
		"jupiter_house_11": "Появляются покровители - ловите возможность. Дружба и социальная активность приносят удачу и возможности. Вы познакомитесь с влиятельными друзьями, расширите свою социальную сеть. Участие в групповых мероприятиях принесет рост и награды. Это хорошее время для **реализации долгосрочных целей**, **получения социальной поддержки**.",
		"jupiter_house_12": "**Духовное расширение** и **внутренние возможности** подчеркнуты. Вы будете **увеличивать духовные возможности**, **расширять внутренний рост**, **духовные поиски** и **внутреннее исследование** выигрывают. Это хорошее время для **духовного развития**, **внутреннего расширения**, **духовных возможностей**. Подходит для **духовного обучения**, **внутреннего исследования**, **духовного роста**.",
		
		// Saturn through houses - complete 12 houses
		"saturn_house_1": "Период усиления **самодисциплины** и **ответственности**. Вы будете **более серьезно относиться к себе и жизненным целям**, можете столкнуться с некоторыми **ограничениями или вызовами**, но все это для **построения более прочной основы**. Это важное время для **воспитания терпения**, **принятия ответственности**, **формирования зрелой личности**.",
		"saturn_house_2": "**Финансовая ответственность** и **материальное управление** становятся фокусом. Вам нужно **серьезно планировать финансы**, **создавать стабильные источники дохода**, можете столкнуться с **финансовыми ограничениями** или необходимостью **сдерживать расходы**. Это хорошее время для **построения долгосрочной финансовой безопасности**, **обучения управлению ценностями**. Подходит для **составления бюджета**, **погашения долгов**, **инвестирования в ценные активы**. Постройте **материальную основу** через **дисциплину и терпение**.",
		"saturn_house_3": "**Ответственность в общении** и **дисциплина в обучении** подчеркнуты. Вам нужно **серьезно относиться к обучению**, **улучшать коммуникативные навыки**, можете почувствовать **ограничения в выражении** или необходимость **говорить более осторожно**. Это хорошее время для **построения базы знаний**, **улучшения отношений с братьями и сестрами**. Подходит для **систематического обучения**, **завершения учебы**, **обработки соседских отношений**. Повысьте **коммуникативные способности** через **терпение и настойчивость**.",
		"saturn_house_4": "**Семейная ответственность** и **эмоциональная зрелость** становятся темой. Вам нужно **принимать семейную ответственность**, **обрабатывать семейные дела**, можете столкнуться с **семейным давлением** или необходимостью **создать семейную структуру**. Это хорошее время для **улучшения семейных отношений**, **обработки недвижимости**, **построения эмоциональной безопасности**. Подходит для **заботы о семье**, **обработки семейного наследия**, **создания семейных традиций**. Постройте **стабильную семейную основу** через **ответственность и обязательства**.",
		"saturn_house_5": "**Творческая ответственность** и **эмоциональная дисциплина** подчеркнуты. Вам нужно **серьезно относиться к творчеству**, **принимать ответственность в отношениях**, можете почувствовать **ограничения в творчестве** или необходимость **более зрело выражать эмоции**. Это хорошее время для **создания творческих привычек**, **обработки отношений с детьми**, **обучения эмоциональным границам**. Подходит для **завершения творческих проектов**, **образования детей**, **создания здоровых привычек развлечений**.",
		"saturn_house_6": "**Ответственность за здоровье** и **рабочая дисциплина** становятся приоритетом. Вам нужно **серьезно относиться к здоровью**, **улучшать рабочие привычки**, можете столкнуться с **вызовами здоровья** или необходимостью **более регулярной жизни**. Это хорошее время для **создания здоровых привычек**, **повышения профессиональных навыков**, **обработки рабочих обязанностей**. Подходит для **медицинских осмотров**, **составления планов здоровья**, **улучшения рабочих процессов**. Постройте **здоровый баланс работы и жизни** через **дисциплину и настойчивость**.",
		"saturn_house_7": "Период испытания отношений и зрелости. Интимные отношения сталкиваются с **реальными испытаниями**, нужно **серьезно относиться к обязательствам и ответственности**. Это может привести к **укреплению отношений** или **очистке нездоровых связей**. Подходит для **создания долгосрочных стабильных партнерских отношений**, **обучения ответственности и границам в отношениях**.",
		"saturn_house_8": "**Глубокая ответственность** и **управление ресурсами** подчеркнуты. Вам нужно **серьезно обрабатывать общие ресурсы**, **сталкиваться с глубокими страхами**, можете столкнуться с **финансовым давлением** или необходимостью **обрабатывать долги**. Это хорошее время для **создания способностей управления ресурсами**, **обработки наследования**, **обучения эмоциональной глубине**. Подходит для **погашения долгов**, **обработки страхования**, **глубокой психологической работы**.",
		"saturn_house_9": "**Ответственность за веру** и **философская дисциплина** становятся темой. Вам нужно **серьезно относиться к вере**, **создавать философскую систему**, можете почувствовать **ограничения в вере** или необходимость **более глубоко учиться**. Это хорошее время для **завершения высшего образования**, **создания системы веры**, **обработки юридических дел**. Подходит для **глубокого обучения**, **обработки юридических документов**, **создания духовных традиций**.",
		"saturn_house_10": "Критический период построения карьеры. Вы столкнетесь с **важными обязанностями и вызовами в карьере**, возможно, потребуется **больше усилий**. Это важный этап для **построения профессиональной репутации**, **реализации долгосрочных карьерных целей**. Успех требует **терпения**, **дисциплины** и **постоянных усилий**.",
		"saturn_house_11": "**Ответственность в дружбе** и **дисциплина в сообществе** подчеркнуты. Вам нужно **серьезно относиться к дружбе**, **принимать ответственность в группах**, можете почувствовать **ограничения в общении** или необходимость **более зрело участвовать в сообществе**. Это хорошее время для **создания долгосрочной дружбы**, **построения положения в группах**, **реализации коллективных целей**. Подходит для **участия в значимых группах**, **создания социальных сетей**, **реализации долгосрочных видений**.",
		"saturn_house_12": "**Внутренняя ответственность** и **духовная дисциплина** становятся темой. Вам нужно **сталкиваться с внутренними страхами**, **обрабатывать карму**, можете почувствовать **необходимость одиночества** или необходимость **более глубоко практиковать**. Это хорошее время для **внутреннего исцеления**, **освобождения от старых паттернов**, **создания духовной дисциплины**. Подходит для **медитации**, **обработки скрытых проблем**, **создания внутренней структуры**.",
		
		// Uranus through houses - complete 12 houses
		"uranus_house_1": "**Личный прорыв** и **внезапные изменения** становятся фокусом. Вы будете **внезапно менять образ**, **прорывать личные ограничения**, **личная свобода** и **независимость** усиливаются. Это хорошее время для **личного прорыва**, **изменения образа**, **погони за свободой**. Подходит для **личных изменений**, **прорыва ограничений**, **погони за независимостью**. Обратите внимание на **стабильность изменений**, избегайте **слишком внезапного**.",
		"uranus_house_2": "**Финансовый прорыв** и **изменение ценностей** подчеркнуты. Вы будете **внезапно менять финансы**, **прорывать ценностные концепции**, **финансовая свобода** и **ценностная независимость** усиливаются. Это хорошее время для **финансового прорыва**, **изменения ценностей**, **погони за финансовой свободой**. Подходит для **финансовых инноваций**, **изменения ценностей**, **финансовой независимости**. Обратите внимание на **финансовую безопасность**, избегайте **слишком рискованного**.",
		"uranus_house_3": "**Прорыв в общении** и **изменение мышления** становятся темой. Вы будете **внезапно менять мышление**, **прорывать способы общения**, **свобода мышления** и **независимость в общении** усиливаются. Это хорошее время для **прорыва мышления**, **изменения общения**, **погони за свободой мышления**. Подходит для **инноваций мышления**, **изменения общения**, **обучения новым технологиям**.",
		"uranus_house_4": "**Семейный прорыв** и **эмоциональные изменения** подчеркнуты. Вы будете **внезапно менять семью**, **прорывать семейные паттерны**, **семейная свобода** и **эмоциональная независимость** усиливаются. Это хорошее время для **семейного прорыва**, **изменения семьи**, **погони за семейной свободой**. Подходит для **семейных изменений**, **обработки недвижимости**, **семейной независимости**. Обратите внимание на **семейную стабильность**, избегайте **слишком внезапного**.",
		"uranus_house_5": "**Творческий прорыв** и **изменение развлечений** становятся фокусом. Вы будете **внезапно менять творчество**, **прорывать способы развлечений**, **творческая свобода** и **независимость в развлечениях** усиливаются. Это хорошее время для **творческого прорыва**, **изменения развлечений**, **погони за творческой свободой**. Подходит для **творческих инноваций**, **изменения развлечений**, **творческой независимости**.",
		"uranus_house_6": "**Рабочий прорыв** и **изменение здоровья** становятся приоритетом. Вы будете **внезапно менять работу**, **прорывать рабочие способы**, **рабочая свобода** и **независимость в здоровье** усиливаются. Это хорошее время для **рабочего прорыва**, **изменения работы**, **погони за рабочей свободой**. Подходит для **рабочих инноваций**, **изменения здоровья**, **рабочей независимости**. Обратите внимание на **рабочую стабильность**, избегайте **слишком внезапного**.",
		"uranus_house_7": "**Прорыв в отношениях** и **изменение сотрудничества** подчеркнуты. Вы будете **внезапно менять отношения**, **прорывать способы сотрудничества**, **свобода в отношениях** и **независимость в сотрудничестве** усиливаются. Это хорошее время для **прорыва в отношениях**, **изменения сотрудничества**, **погони за свободой в отношениях**. Подходит для **инноваций в отношениях**, **изменения сотрудничества**, **независимости в отношениях**. Обратите внимание на **стабильность отношений**, избегайте **слишком внезапного**.",
		"uranus_house_8": "**Прорыв ресурсов** и **изменение глубины** становятся темой. Вы будете **внезапно менять ресурсы**, **прорывать глубокие отношения**, **свобода ресурсов** и **независимость в глубине** усиливаются. Это хорошее время для **прорыва ресурсов**, **изменения глубины**, **погони за свободой ресурсов**. Подходит для **инноваций ресурсов**, **изменения глубины**, **независимости ресурсов**. Обратите внимание на **безопасность ресурсов**, избегайте **слишком рискованного**.",
		"uranus_house_9": "**Духовный прорыв** и **изменение идеалов** подчеркнуты. Вы будете **внезапно менять идеалы**, **прорывать духовные способы**, **духовная свобода** и **независимость в идеалах** усиливаются. Это хорошее время для **духовного прорыва**, **изменения идеалов**, **погони за духовной свободой**. Подходит для **духовных инноваций**, **изменения идеалов**, **духовной независимости**.",
		"uranus_house_10": "**Карьерный прорыв** и **изменение карьеры** становятся фокусом. Вы будете **внезапно менять карьеру**, **прорывать карьерные способы**, **карьерная свобода** и **независимость в карьере** усиливаются. Это хорошее время для **карьерного прорыва**, **изменения карьеры**, **погони за карьерной свободой**. Подходит для **карьерных инноваций**, **изменения карьеры**, **карьерной независимости**. Обратите внимание на **карьерную стабильность**, избегайте **слишком внезапного**.",
		"uranus_house_11": "**Прорыв сообщества** и **изменение идеалов** подчеркнуты. Вы будете **внезапно менять сообщество**, **прорывать идеальные способы**, **свобода сообщества** и **независимость в идеалах** усиливаются. Это хорошее время для **прорыва сообщества**, **изменения идеалов**, **погони за свободой сообщества**. Подходит для **инноваций сообщества**, **изменения идеалов**, **независимости сообщества**.",
		"uranus_house_12": "**Внутренний прорыв** и **духовные изменения** становятся темой. Вы будете **внезапно менять внутреннее**, **прорывать духовные способы**, **внутренняя свобода** и **духовная независимость** усиливаются. Это хорошее время для **внутреннего прорыва**, **изменения духовного**, **погони за внутренней свободой**. Подходит для **внутренних инноваций**, **изменения духовного**, **внутренней независимости**.",
		
		// Neptune through houses - complete 12 houses
		"neptune_house_1": "**Личный идеал** и **духовное вдохновение** подчеркнуты. Вы будете **идеализировать себя**, **стремиться к духовному вдохновению**, **личный идеал** и **духовное вдохновение** увеличиваются. Это хорошее время для **личного идеала**, **духовного вдохновения**, **погони за идеальным собой**. Подходит для **идеального себя**, **духовного вдохновения**, **личного идеала**. Обратите внимание на **реальность личности**, избегайте **чрезмерной идеализации**.",
		"neptune_house_2": "**Финансовый идеал** и **ценностное вдохновение** становятся темой. Вы будете **идеализировать финансы**, **стремиться к ценностному вдохновению**, **финансовый идеал** и **ценностное вдохновение** увеличиваются. Это хорошее время для **финансового идеала**, **ценностного вдохновения**, **погони за идеальными финансами**. Подходит для **идеальных финансов**, **ценностного вдохновения**, **финансового идеала**. Обратите внимание на **финансовую реальность**, избегайте **чрезмерной идеализации**.",
		"neptune_house_3": "**Идеал общения** и **мыслительное вдохновение** подчеркнуты. Вы будете **идеализировать общение**, **стремиться к мыслительному вдохновению**, **идеал общения** и **мыслительное вдохновение** увеличиваются. Это хорошее время для **идеала общения**, **мыслительного вдохновения**, **погони за идеальным общением**. Подходит для **идеального общения**, **мыслительного вдохновения**, **идеала общения**.",
		"neptune_house_4": "**Семейный идеал** и **эмоциональное вдохновение** становятся темой. Вы будете **идеализировать семью**, **стремиться к эмоциональному вдохновению**, **семейный идеал** и **эмоциональное вдохновение** увеличиваются. Это хорошее время для **семейного идеала**, **эмоционального вдохновения**, **погони за идеальной семьей**. Подходит для **идеальной семьи**, **эмоционального вдохновения**, **семейного идеала**. Обратите внимание на **семейную реальность**, избегайте **чрезмерной идеализации**.",
		"neptune_house_5": "**Творческий идеал** и **развлекательное вдохновение** подчеркнуты. Вы будете **идеализировать творчество**, **стремиться к развлекательному вдохновению**, **творческий идеал** и **развлекательное вдохновение** увеличиваются. Это хорошее время для **творческого идеала**, **развлекательного вдохновения**, **погони за идеальным творчеством**. Подходит для **идеального творчества**, **развлекательного вдохновения**, **творческого идеала**.",
		"neptune_house_6": "**Рабочий идеал** и **вдохновение здоровья** становятся приоритетом. Вы будете **идеализировать работу**, **стремиться к вдохновению здоровья**, **рабочий идеал** и **вдохновение здоровья** увеличиваются. Это хорошее время для **рабочего идеала**, **вдохновения здоровья**, **погони за идеальной работой**. Подходит для **идеальной работы**, **вдохновения здоровья**, **рабочего идеала**. Обратите внимание на **рабочую реальность**, избегайте **чрезмерной идеализации**.",
		"neptune_house_7": "**Идеал отношений** и **вдохновение сотрудничества** подчеркнуты. Вы будете **идеализировать отношения**, **стремиться к вдохновению сотрудничества**, **идеал отношений** и **вдохновение сотрудничества** увеличиваются. Это хорошее время для **идеала отношений**, **вдохновения сотрудничества**, **погони за идеальными отношениями**. Подходит для **идеальных отношений**, **вдохновения сотрудничества**, **идеала отношений**. Обратите внимание на **реальность отношений**, избегайте **чрезмерной идеализации**.",
		"neptune_house_8": "**Идеал ресурсов** и **глубокое вдохновение** становятся темой. Вы будете **идеализировать ресурсы**, **стремиться к глубокому вдохновению**, **идеал ресурсов** и **глубокое вдохновение** увеличиваются. Это хорошее время для **идеала ресурсов**, **глубокого вдохновения**, **погони за идеальными ресурсами**. Подходит для **идеальных ресурсов**, **глубокого вдохновения**, **идеала ресурсов**. Обратите внимание на **реальность ресурсов**, избегайте **чрезмерной идеализации**.",
		"neptune_house_9": "**Духовный идеал** и **идеальное вдохновение** подчеркнуты. Вы будете **идеализировать духовное**, **стремиться к идеальному вдохновению**, **духовный идеал** и **идеальное вдохновение** увеличиваются. Это хорошее время для **духовного идеала**, **идеального вдохновения**, **погони за идеальным духовным**. Подходит для **идеального духовного**, **идеального вдохновения**, **духовного идеала**.",
		"neptune_house_10": "**Карьерный идеал** и **карьерное вдохновение** подчеркнуты. Вы будете **идеализировать карьеру**, **стремиться к карьерному вдохновению**, **карьерный идеал** и **карьерное вдохновение** увеличиваются. Это хорошее время для **карьерного идеала**, **карьерного вдохновения**, **погони за идеальной карьерой**. Подходит для **идеальной карьеры**, **карьерного вдохновения**, **карьерного идеала**. Обратите внимание на **карьерную реальность**, избегайте **чрезмерной идеализации**.",
		"neptune_house_11": "**Идеал сообщества** и **коллективное вдохновение** подчеркнуты. Вы будете **идеализировать сообщество**, **стремиться к коллективному вдохновению**, **идеал сообщества** и **коллективное вдохновение** увеличиваются. Это хорошее время для **идеала сообщества**, **коллективного вдохновения**, **погони за идеальным сообществом**. Подходит для **идеального сообщества**, **коллективного вдохновения**, **идеала сообщества**.",
		"neptune_house_12": "**Внутренний идеал** и **духовное вдохновение** становятся темой. Вы будете **идеализировать внутреннее**, **стремиться к духовному вдохновению**, **внутренний идеал** и **духовное вдохновение** увеличиваются. Это хорошее время для **внутреннего идеала**, **духовного вдохновения**, **погони за идеальным внутренним**. Подходит для **идеального внутреннего**, **духовного вдохновения**, **внутреннего идеала**.",
		
		// Pluto through houses - complete 12 houses
		"pluto_house_1": "**Личная трансформация** и **глубокие изменения** подчеркнуты. Вы будете **испытывать глубокую трансформацию**, **менять личную сущность**, **личная сила** и **способность к трансформации** увеличиваются. Это хорошее время для **личной трансформации**, **глубоких изменений**, **погони за личной силой**. Подходит для **личной трансформации**, **глубоких изменений**, **личной силы**. Обратите внимание на **стабильность трансформации**, избегайте **слишком интенсивного**.",
		"pluto_house_2": "**Финансовая трансформация** и **изменение ценностей** подчеркнуты. Вы будете **испытывать финансовую трансформацию**, **менять ценностные концепции**, **финансовая сила** и **ценностная трансформация** увеличиваются. Это хорошее время для **финансовой трансформации**, **изменения ценностей**, **погони за финансовой силой**. Подходит для **финансовой трансформации**, **изменения ценностей**, **финансовой силы**. Обратите внимание на **финансовую безопасность**, избегайте **слишком интенсивного**.",
		"pluto_house_3": "**Трансформация общения** и **изменение мышления** становятся темой. Вы будете **испытывать трансформацию общения**, **менять способы мышления**, **сила общения** и **трансформация мышления** увеличиваются. Это хорошее время для **трансформации общения**, **изменения мышления**, **погони за силой общения**. Подходит для **трансформации общения**, **изменения мышления**, **силы общения**.",
		"pluto_house_4": "**Семейная трансформация** и **эмоциональные изменения** подчеркнуты. Вы будете **испытывать семейную трансформацию**, **менять семейные паттерны**, **семейная сила** и **эмоциональная трансформация** увеличиваются. Это хорошее время для **семейной трансформации**, **эмоциональных изменений**, **погони за семейной силой**. Подходит для **семейной трансформации**, **эмоциональных изменений**, **семейной силы**. Обратите внимание на **семейную стабильность**, избегайте **слишком интенсивного**.",
		"pluto_house_5": "**Творческая трансформация** и **изменение развлечений** подчеркнуты. Вы будете **испытывать творческую трансформацию**, **менять способы развлечений**, **творческая сила** и **трансформация развлечений** увеличиваются. Это хорошее время для **творческой трансформации**, **изменения развлечений**, **погони за творческой силой**. Подходит для **творческой трансформации**, **изменения развлечений**, **творческой силы**.",
		"pluto_house_6": "**Рабочая трансформация** и **изменение здоровья** становятся приоритетом. Вы будете **испытывать рабочую трансформацию**, **менять рабочие способы**, **рабочая сила** и **трансформация здоровья** увеличиваются. Это хорошее время для **рабочей трансформации**, **изменения здоровья**, **погони за рабочей силой**. Подходит для **рабочей трансформации**, **изменения здоровья**, **рабочей силы**. Обратите внимание на **рабочую стабильность**, избегайте **слишком интенсивного**.",
		"pluto_house_7": "**Трансформация отношений** и **изменение сотрудничества** подчеркнуты. Вы будете **испытывать трансформацию отношений**, **менять способы сотрудничества**, **сила отношений** и **трансформация сотрудничества** увеличиваются. Это хорошее время для **трансформации отношений**, **изменения сотрудничества**, **погони за силой отношений**. Подходит для **трансформации отношений**, **изменения сотрудничества**, **силы отношений**. Обратите внимание на **стабильность отношений**, избегайте **слишком интенсивного**.",
		"pluto_house_8": "**Трансформация ресурсов** и **глубокие изменения** становятся темой. Вы будете **испытывать трансформацию ресурсов**, **менять глубокие отношения**, **сила ресурсов** и **глубокая трансформация** увеличиваются. Это хорошее время для **трансформации ресурсов**, **глубоких изменений**, **погони за силой ресурсов**. Подходит для **трансформации ресурсов**, **изменения инвестиций**, **силы ресурсов**. Обратите внимание на **безопасность ресурсов**, избегайте **слишком интенсивного**.",
		"pluto_house_9": "**Духовная трансформация** и **изменение идеалов** подчеркнуты. Вы будете **испытывать духовную трансформацию**, **менять идеальные способы**, **духовная сила** и **трансформация идеалов** увеличиваются. Это хорошее время для **духовной трансформации**, **изменения идеалов**, **погони за духовной силой**. Подходит для **духовной трансформации**, **изменения идеалов**, **духовной силы**.",
		"pluto_house_10": "**Карьерная трансформация** и **изменение карьеры** подчеркнуты. Вы будете **испытывать карьерную трансформацию**, **менять карьерные способы**, **карьерная сила** и **трансформация карьеры** увеличиваются. Это хорошее время для **карьерной трансформации**, **изменения карьеры**, **погони за карьерной силой**. Подходит для **карьерной трансформации**, **изменения карьеры**, **карьерной силы**. Обратите внимание на **карьерную стабильность**, избегайте **слишком интенсивного**.",
		"pluto_house_11": "**Трансформация сообщества** и **изменение идеалов** подчеркнуты. Вы будете **испытывать трансформацию сообщества**, **менять идеальные способы**, **сила сообщества** и **трансформация идеалов** увеличиваются. Это хорошее время для **трансформации сообщества**, **изменения идеалов**, **погони за силой сообщества**. Подходит для **трансформации сообщества**, **изменения идеалов**, **силы сообщества**.",
		"pluto_house_12": "**Внутренняя трансформация** и **духовные изменения** становятся темой. Вы будете **испытывать внутреннюю трансформацию**, **менять духовные способы**, **внутренняя сила** и **духовная трансформация** увеличиваются. Это хорошее время для **внутренней трансформации**, **духовных изменений**, **погони за внутренней силой**. Подходит для **внутренней трансформации**, **духовных изменений**, **внутренней силы**.",
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
	// Fallback: reuse aspect interpretation with progression context
	if aspectText := getChineseAspectInterpretation(key, isPositive); aspectText != "" {
		return "推运层面，这一相位反映长期主题：" + aspectText
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
	// Fallback: reuse aspect interpretation with progression context
	if aspectText := getEnglishAspectInterpretation(key, isPositive); aspectText != "" {
		return "As a progression, this aspect reflects a longer-term theme: " + aspectText
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
	// Fallback: reuse aspect interpretation with progression context
	if aspectText := getRussianAspectInterpretation(key, isPositive); aspectText != "" {
		return "В прогрессии этот аспект отражает долгосрочную тему: " + aspectText
	}
	if isPositive {
		return "Это период позитивного развития и личностного роста. Вы испытаете гармоничный прогресс в этой области с усилением внутренних способностей."
	}
	return "Это период столкновения с вызовами и содействия росту. Вы столкнетесь с некоторыми трудностями в этой области, но эти переживания принесут глубокую зрелость и понимание."
}

// ========== Aspect Interpretations ==========
// getChineseAspectInterpretation is in aspect_interpretations_zh.go
// getEnglishAspectInterpretation is in aspect_interpretations_en.go
// getRussianAspectInterpretation is in aspect_interpretations_ru.go

// ========== Retrograde Interpretations ==========

func getChineseRetrogradeInterpretation(planet models.PlanetID) string {
	texts := map[string]string{
		"mercury": "**水星逆行**期间，**沟通、合约、出行、电子设备**容易有延误或反复。适合**复盘旧项目**、**整理文档**、**修复关系**，而非开启全新合作或签大单。**工作沟通**请多确认、留书面记录；**学业**上宜复习多于突击。保持耐心与弹性，把逆行当作**内省与修正**的窗口。",
		"venus":   "**金星逆行**时，**感情、审美、金钱、人际**易有回顾与调整。适合**厘清真正想要的关系与价值**，**修补旧怨**或**重审消费习惯**，不宜草率表白、整容或做大额投资。**关系**上多倾听、少计较；**财务**上以守成为主。利用这段时间**自爱**与**价值澄清**。",
		"mars":    "**火星逆行**期间，**行动力、竞争、冲突、体能**会显得反复或受阻。适合**重新规划目标**、**处理积压事项**、**运动康复**，避免冲动决策、与人硬刚或开启新诉讼。**事业**上以稳为主；**健康**注意炎症与外伤。把能量用在**策略与复盘**上。",
		"jupiter": "**木星逆行**时，**扩张、机会、信念、远行**会放缓或需要重新评估。适合**检视现有计划是否过大**、**修正过度乐观**、**内化信仰与哲学**，不宜盲目扩张、乱投资或轻信承诺。**学业与法律**事宜宜谨慎推进。借逆行**夯实基础**而非追逐新机会。",
		"saturn":  "**土星逆行**期间，**责任、结构、权威、时间**主题被强化。适合**面对未完成的责任**、**调整长期规划**、**与长辈或上司理顺关系**，可能感到**被拖延或批评**。**事业与关系**上的考验是常态；用**耐心与纪律**一步步化解，而非逃避。",
		"uranus":  "**天王星逆行**时，**变革、自由、突发、创新**更多转向内在或延迟表现。适合**审视自己的叛逆与固执**、**调整生活方式**而非强求外界剧变，**科技与社群**相关计划可能反复。保持**内在弹性**，等顺行后再大力推动改革。",
		"neptune": "**海王星逆行**期间，**灵感、慈悲、边界、逃避**被检视。适合**分清理想与幻想**、**处理成瘾或依赖**、**艺术与灵性**上的沉淀，不宜在这时做重大牺牲或签模糊合约。**关系**中易有误会；**财务**需防骗。用逆行**落地梦想**、**设好边界**。",
		"pluto":   "**冥王星逆行**时，**权力、执念、转化、阴影**更多在内心发酵。适合**面对恐惧与执著**、**结束不再服务你的模式**、**资源与关系的深度清理**，不宜强行控制他人或做毁灭式决定。**危机感**是转化前兆；用**觉察与放下**代替对抗。",
	}
	if s, ok := texts[string(planet)]; ok {
		return s
	}
	return "该行星逆行期间，相关生活领域宜**回顾与调整**，少做全新开端，多**修正与内化**。"
}

func getEnglishRetrogradeInterpretation(planet models.PlanetID) string {
	texts := map[string]string{
		"mercury": "During **Mercury retrograde**, **communication, contracts, travel, and devices** are prone to delays or mix-ups. Favor **reviewing old projects**, **organizing paperwork**, and **mending ties** over new deals or big commitments. **Confirm in writing** at work; **study** by revising rather than cramming. Use this period for **reflection and correction**.",
		"venus":   "When **Venus is retrograde**, **love, beauty, money, and relationships** go through review. Good for **clarifying what you truly value**, **healing old rifts**, or **revisiting spending**; avoid rushed confessions, major cosmetic changes, or large investments. **Listen more** in relationships; **preserve** financially. Use the time for **self-worth and clarity**.",
		"mars":    "During **Mars retrograde**, **drive, competition, conflict, and stamina** can feel blocked or uneven. Favor **replanning goals**, **clearing backlog**, and **recovery** over new fights or lawsuits. **Career**: stay steady; **health**: watch inflammation and injury. Channel energy into **strategy and review**.",
		"jupiter": "When **Jupiter is retrograde**, **expansion, opportunity, belief, and travel** slow or need reassessment. Good for **checking if plans are overstretched**, **tempering optimism**, and **internalizing philosophy**; avoid over-expanding, speculative bets, or trusting big promises. **Education and legal** matters: proceed with care. **Consolidate** rather than chase new luck.",
		"saturn":  "During **Saturn retrograde**, **duty, structure, authority, and time** are emphasized. Good for **facing unfinished responsibilities**, **adjusting long-term plans**, and **sorting dynamics with elders or bosses**. You may feel **delayed or criticized**. **Career and relationships** are under review; meet them with **patience and discipline**.",
		"uranus":  "When **Uranus is retrograde**, **change, freedom, shocks, and innovation** turn inward or delay. Good for **examining your own rebellion and rigidity**, **adjusting lifestyle** rather than forcing outer upheaval. **Tech and community** plans may shift. Stay **internally flexible**; push big changes after it goes direct.",
		"neptune": "During **Neptune retrograde**, **inspiration, compassion, boundaries, and escape** are under review. Good for **separating dream from delusion**, **addressing addiction or dependency**, and **digesting art and spirit**; avoid major sacrifices or vague contracts. **Relationships**: misunderstandings likely; **finance**: guard against deception. Use it to **ground dreams** and **set boundaries**.",
		"pluto":   "When **Pluto is retrograde**, **power, obsession, transformation, and shadow** work inwardly. Good for **facing fear and attachment**, **ending patterns that no longer serve**, and **deep clearing of resources and ties**; avoid controlling others or destructive decisions. **Crisis feeling** can precede renewal; **awareness and release** over fight.",
	}
	if s, ok := texts[string(planet)]; ok {
		return s
	}
	return "During this planet's retrograde, related areas of life benefit from **review and adjustment**; favor **revision and integration** over brand-new starts."
}

func getRussianRetrogradeInterpretation(planet models.PlanetID) string {
	texts := map[string]string{
		"mercury": "В период **ретроградного Меркурия** **общение, контракты, поездки и техника** склонны к задержкам и путанице. Уместны **разбор старых проектов**, **упорядочивание документов**, **восстановление связей** — лучше отложить новые сделки. **На работе** фиксируйте договорённости; **учёба**: повторение важнее авралов. Используйте время для **внутренней проверки и исправлений**.",
		"venus":   "При **ретроградной Венере** **любовь, красота, деньги и отношения** пересматриваются. Хорошо **прояснить истинные ценности**, **залечить старые обиды**, **пересмотреть траты**; не спешите с признаниями, крупными изменениями во внешности или вложениями. **В отношениях** больше слушайте; **финансы** берегите. Время для **самооценки и ясности**.",
		"mars":    "В **ретроградный Марс** **драйв, конкуренция, конфликты и выносливость** могут блокироваться. Уместны **перепланирование целей**, **разбор завалов**, **восстановление**; избегайте новых конфликтов и судов. **Карьера**: стабильность; **здоровье**: внимание к воспалениям и травмам. Направляйте энергию в **стратегию и анализ**.",
		"jupiter": "При **ретроградном Юпитере** **расширение, возможности, вера и поездки** замедляются. Хорошо **проверить, не переоценены ли планы**, **снизить избыточный оптимизм**, **углубить философию**; не расширяйтесь вслепую и не доверяйте громким обещаниям. **Учёба и юридическое** — осторожно. **Укрепляйте основу**, а не гонитесь за новым.",
		"saturn":  "В период **ретроградного Сатурна** усилены темы **долга, структуры, авторитета и времени**. Уместны **разбор незавершённых обязанностей**, **коррекция долгосрочных планов**, **выстраивание отношений с начальством и старшими**. Возможны **задержки и критика**. **Карьера и отношения** под проверкой; отвечайте **терпением и дисциплиной**.",
		"uranus":  "При **ретроградном Уране** **перемены, свобода, неожиданности и инновации** уходят внутрь или откладываются. Хорошо **посмотреть на свой бунт и ригидность**, **скорректировать образ жизни**, а не провоцировать внешний переворот. **Технологии и сообщества** могут меняться. Сохраняйте **внутреннюю гибкость**; крупные реформы — после директа.",
		"neptune": "В **ретроградный Нептун** пересматриваются **вдохновение, сострадание, границы и побег**. Хорошо **отделять мечту от иллюзии**, **работать с зависимостями**, **переваривать искусство и духовность**; избегайте жертв и размытых договоров. **Отношения**: недопонимания; **финансы**: остерегайтесь обмана. Используйте время, чтобы **приземлить мечты** и **обозначить границы**.",
		"pluto":   "При **ретроградном Плутоне** **власть, одержимость, трансформация и тень** работают изнутри. Хорошо **встречать страхи и привязанности**, **завершать отжившие сценарии**, **глубоко чистить ресурсы и связи**; не контролируйте других и не принимайте разрушительных решений. **Чувство кризиса** может предшествовать обновлению; **осознанность и отпускание** вместо борьбы.",
	}
	if s, ok := texts[string(planet)]; ok {
		return s
	}
	return "В ретроградный период этой планеты связанные сферы жизни лучше **пересматривать и корректировать**; предпочтите **доработку и интеграцию** новым стартам."
}

// ========== Profection Lord (Annual Lord) Interpretations ==========

func getChineseProfectionLordInterpretation(planet models.PlanetID) string {
	texts := map[string]string{
		"sun":     "今年**年主星是太阳**，整年的主题围绕**自我、身份、活力与领导力**。**健康**与**事业**是重点：适合确立**个人目标**、**展现专业能力**、**照顾身体**。你会更在意**被看见与被认可**；与**父亲、权威、上司**的关系也会被强调。利用这一年**夯实自我认同**、**在职场或生活中担当主角**。",
		"moon":    "今年**年主星是月亮**，主题落在**情绪、家庭、安全感与内在需求**。**家庭关系**、**居住环境**、**与母亲或女性的关系**以及**身心健康**会被放大。适合**经营家庭**、**处理房产**、**照顾情绪**；**事业**上也可能与**公众、照顾、餐饮**相关。这一年**直觉**增强，适合**向内扎根**、**建立情感安全**。",
		"mercury": "今年**年主星是水星**，**沟通、学习、短途出行、兄弟邻里**成为年度重点。**学业、写作、签约、传播**易有进展；**职场沟通**与**人脉**也重要。适合**考学、培训、谈判**，注意**信息准确**与**合同条款**。这一年思维活跃，宜**多学多表达**、**连接信息与人**。",
		"venus":   "今年**年主星是金星**，**爱情、金钱、审美、人际**是核心。**感情**与**合作**有机会推进；**财务**与**享受**被强调。适合**深化关系**、**理财规划**、**艺术或美容**相关事务；注意**价值取舍**与**关系边界**。这一年**魅力**与**资源**易被看见，用好**爱与价值**的主题。",
		"mars":    "今年**年主星是火星**，**行动、竞争、勇气、体能**被激活。**事业**上易有**冲刺与竞争**；**健康**需留意**发炎、外伤、过劳**。适合**设定目标并执行**、**运动健身**、**解决冲突**；避免**冲动与树敌**。这一年**斗志**强，宜**把能量用在建设性目标**上。",
		"jupiter": "今年**年主星是木星**，**扩展、机会、信念、远行**成为年度基调。**学业、法律、出版、跨地域**事务易有进展；**贵人**与**资源**可能增多。适合**设定更大目标**、**投资成长**、**探索信仰与哲学**；避免**过度乐观与铺张**。这一年**机遇**多，宜**顺势而为、稳健扩张**。",
		"saturn":  "今年**年主星是土星**，**责任、结构、时间、权威**是主题。**事业**与**关系**可能面临**考验与成熟**；**健康**需**规律与纪律**。适合**兑现承诺**、**建立长期结构**、**与长辈或规则共处**；避免**逃避责任**。这一年**坚持**会带来**认可**，宜**一步一个脚印**。",
	}
	if s, ok := texts[string(planet)]; ok {
		return s
	}
	return "今年年主星影响整年主题，相关生活领域会得到**持续关注**；宜根据该行星象征的领域**规划与落实**。"
}

func getEnglishProfectionLordInterpretation(planet models.PlanetID) string {
	texts := map[string]string{
		"sun":     "This year **the annual lord is the Sun**; themes center on **self, identity, vitality, and leadership**. **Health** and **career** are in focus: good for **setting personal goals**, **showing competence**, **self-care**. You may care more about **being seen and recognized**; **father, authority, superiors** are highlighted. Use the year to **solidify self-identity** and **take the lead at work or in life**.",
		"moon":    "This year **the annual lord is the Moon**; themes are **emotions, family, security, and inner needs**. **Family, home, relationship with mother or women**, and **wellbeing** are emphasized. Good for **nurturing home**, **property**, **emotional care**; **career** may involve **public, care, or food**. **Intuition** is stronger; use the year to **root inwardly** and **build emotional safety**.",
		"mercury": "This year **the annual lord is Mercury**; **communication, learning, short travel, siblings and neighbors** are in focus. **Study, writing, contracts, media** can progress; **workplace communication** and **networks** matter. Good for **exams, training, negotiation**; mind **accuracy and contracts**. The year is mentally active — **learn, express, connect**.",
		"venus":   "This year **the annual lord is Venus**; **love, money, beauty, relationships** are central. **Romance** and **partnership** can advance; **finance** and **pleasure** are emphasized. Good for **deepening relationships**, **financial planning**, **arts or beauty**; watch **values and boundaries**. **Charm** and **resources** are visible — lean into **love and value**.",
		"mars":    "This year **the annual lord is Mars**; **action, competition, courage, stamina** are active. **Career** may see **drive and rivalry**; **health**: watch **inflammation, injury, overwork**. Good for **goals and execution**, **exercise**, **resolving conflict**; avoid **impulse and enemies**. **Channel energy** into constructive aims.",
		"jupiter": "This year **the annual lord is Jupiter**; **expansion, opportunity, belief, travel** set the tone. **Education, law, publishing, abroad** can progress; **benefactors** and **resources** may increase. Good for **larger goals**, **investing in growth**, **exploring philosophy**; avoid **over-optimism and excess**. **Opportunity** is high — **expand with care**.",
		"saturn":  "This year **the annual lord is Saturn**; **responsibility, structure, time, authority** are themes. **Career** and **relationships** may face **tests and maturity**; **health** benefits from **routine and discipline**. Good for **keeping commitments**, **building long-term structure**, **working with elders and rules**; avoid **evading duty**. **Persistence** brings **recognition** — **step by step**.",
	}
	if s, ok := texts[string(planet)]; ok {
		return s
	}
	return "The annual lord shapes the year's themes; related areas receive **ongoing emphasis** — **plan and act** according to that planet's symbolism."
}

func getRussianProfectionLordInterpretation(planet models.PlanetID) string {
	texts := map[string]string{
		"sun":     "В этом году **годовым управителем является Солнце**; темы — **самость, идентичность, жизненная сила и лидерство**. В фокусе **здоровье** и **карьера**: хорошее время для **личных целей**, **проявления компетентности**, **заботы о себе**. Важнее **быть замеченным**; актуальны **отец, авторитет, начальство**. Используйте год для **укрепления самоидентичности** и **роли лидера**.",
		"moon":    "В этом году **годовым управителем является Луна**; темы — **эмоции, семья, безопасность, внутренние потребности**. Усилены **семья, дом, отношения с матерью или женщинами**, **благополучие**. Хорошо для **домашнего уюта**, **недвижимости**, **эмоциональной заботы**; **карьера** может быть связана с **публикой, заботой, питанием**. **Интуиция** сильнее; **углубляйтесь внутрь** и **стройте эмоциональную безопасность**.",
		"mercury": "В этом году **годовым управителем является Меркурий**; в фокусе **общение, учёба, короткие поездки, братья и соседи**. **Учёба, письмо, контракты, медиа** могут продвинуться; важны **коммуникация на работе** и **сети**. Хорошо для **экзаменов, тренингов, переговоров**; следите за **точностью и договорами**. Год активного мышления — **учитесь, выражайтесь, связывайте**.",
		"venus":   "В этом году **годовым управителем является Венера**; в центре **любовь, деньги, красота, отношения**. **Романтика** и **партнёрство** могут развиваться; усилены **финансы** и **удовольствие**. Хорошо для **углубления отношений**, **финансового планирования**, **искусства и красоты**; следите за **ценностями и границами**. **Очарование** и **ресурсы** на виду — опора на **любовь и ценность**.",
		"mars":    "В этом году **годовым управителем является Марс**; активны **действие, соревнование, смелость, выносливость**. **Карьера** может дать **рывок и конкуренцию**; **здоровье**: внимание **воспалению, травмам, переутомлению**. Хорошо для **целей и исполнения**, **спорта**, **разрешения конфликтов**; избегайте **импульса и врагов**. Направляйте энергию в **конструктивные цели**.",
		"jupiter": "В этом году **годовым управителем является Юпитер**; тон задают **расширение, возможности, вера, путешествия**. **Образование, право, издательства, зарубеж** могут продвинуться; **покровители** и **ресурсы** могут возрасти. Хорошо для **крупных целей**, **инвестиций в рост**, **философии**; избегайте **переоптимизма и избытка**. **Возможностей** много — **расширяйтесь с умом**.",
		"saturn":  "В этом году **годовым управителем является Сатурн**; темы — **ответственность, структура, время, авторитет**. **Карьера** и **отношения** могут пройти **проверку и взросление**; **здоровье** выигрывает от **режима и дисциплины**. Хорошо для **обязательств**, **долгосрочной структуры**, **отношений со старшими и правилами**; не уходите от **долга**. **Настойчивость** приносит **признание** — **шаг за шагом**.",
	}
	if s, ok := texts[string(planet)]; ok {
		return s
	}
	return "Годовой управитель задаёт темы года; связанные сферы получают **постоянное внимание** — **планируйте и действуйте** в соответствии с символикой планеты."
}

// ========== Planetary Hour Interpretations (hourly) ==========

func getChinesePlanetaryHourInterpretation(planet models.PlanetID) string {
	texts := map[string]string{
		"sun":     "当前处于**太阳时**，能量偏向**自我表达、权威与活力**。适合**拍板决定**、**展现领导力**、**处理与父亲或上司相关**的事；**重要会面、签约、健康相关**也较顺。避免在此时**过度妥协或隐藏主见**。",
		"moon":    "当前处于**月亮时**，氛围偏**情绪与直觉**。适合**处理家事、照顾他人、饮食与休息**；**情感沟通、内心复盘**较顺。重大**签约或新项目启动**可考虑延后；宜**倾听感受、巩固安全感**。",
		"mercury": "当前处于**水星时**，利于**沟通、学习、短途与文书**。适合**开会、邮件、签约、考试、短途出行**；**谈判与信息交换**效率高。避免此时做**重大情感或长期投资**决策；宜**多表达、多记录**。",
		"venus":   "当前处于**金星时**，适合**关系、金钱与审美**相关事务。**约会、合作、理财、美容、艺术**较顺；**表达好感、谈条件**易被接受。不宜**冲动冲突或重大冒险**；宜**维系和谐、做让自己愉悦的事**。",
		"mars":    "当前处于**火星时**，能量偏**行动与竞争**。适合**运动、截止日前冲刺、解决冲突、争取权益**；**启动项目、面试、竞赛**可借势。避免**无谓争吵或冒险**；宜**目标明确、执行果断**。",
		"jupiter": "当前处于**木星时**，利于**扩展、机会与信念**。**求学、法律、出版、远行、大额规划**较顺；**求助贵人、拓展人脉**易有回应。不宜**过度承诺或挥霍**；宜**顺势而为、留有余地**。",
		"saturn":  "当前处于**土星时**，适合**责任、结构与长期事务**。**汇报、承诺、纪律、与长辈/规则打交道**较顺；**收尾、复盘、制定计划**效率高。不宜**冒险或轻率承诺**；宜**务实、一步一个脚印**。",
	}
	if s, ok := texts[string(planet)]; ok {
		return s
	}
	return "当前行星时强调该行星所象征的领域；可据此**安排相应类型的活动**，事半功倍。"
}

func getEnglishPlanetaryHourInterpretation(planet models.PlanetID) string {
	texts := map[string]string{
		"sun":     "You're in a **Sun hour** — energy favors **self-expression, authority, and vitality**. Good for **decisions**, **showing leadership**, **matters with father or superiors**; **meetings, contracts, health** also flow. Avoid **over-compromising or hiding your stance** now.",
		"moon":    "You're in a **Moon hour** — mood is **emotional and intuitive**. Good for **home, caring for others, food and rest**; **heart-to-heart talks, inner review** work well. Consider postponing **big contracts or new launches**; **listen to feelings, build security**.",
		"mercury": "You're in a **Mercury hour** — ideal for **communication, learning, short trips, paperwork**. Good for **meetings, emails, contracts, exams, short travel**; **negotiation and info exchange** are efficient. Avoid **major emotional or long-term money decisions**; **speak and write**.",
		"venus":   "You're in a **Venus hour** — good for **relationships, money, and beauty**. **Dating, partnership, finance, looks, arts** flow; **expressing affection, discussing terms** are well received. Avoid **confrontation or big risks**; **keep harmony, do what pleases you**.",
		"mars":    "You're in a **Mars hour** — energy is **action and competition**. Good for **exercise, deadline pushes, resolving conflict, standing your ground**; **launching projects, interviews, competition** can use this. Avoid **pointless arguments or recklessness**; **be clear and decisive**.",
		"jupiter": "You're in a **Jupiter hour** — favors **expansion, opportunity, and belief**. **Study, law, publishing, travel, big plans** flow; **asking for help, networking** get responses. Avoid **over-committing or overspending**; **go with the flow, leave room**.",
		"saturn":  "You're in a **Saturn hour** — good for **duty, structure, and long-term matters**. **Reporting, commitments, discipline, dealing with elders or rules** flow; **wrapping up, review, planning** are efficient. Avoid **gambles or rash promises**; **be practical, step by step**.",
	}
	if s, ok := texts[string(planet)]; ok {
		return s
	}
	return "The current planetary hour emphasizes that planet's domains; **schedule related activities** for better results."
}

func getRussianPlanetaryHourInterpretation(planet models.PlanetID) string {
	texts := map[string]string{
		"sun":     "Сейчас **час Солнца** — энергия в пользу **самовыражения, авторитета и жизненной силы**. Хорошо для **решений**, **лидерства**, **дел с отцом или начальством**; **встречи, контракты, здоровье** тоже в потоке. Избегайте **чрезмерных уступок или сокрытия позиции**.",
		"moon":    "Сейчас **час Луны** — настроение **эмоциональное и интуитивное**. Хорошо для **дома, заботы о других, еды и отдыха**; **сердечные разговоры, внутренний разбор** уместны. Крупные **контракты или новые старты** лучше перенести; **слушайте чувства, укрепляйте безопасность**.",
		"mercury": "Сейчас **час Меркурия** — идеально для **общения, учёбы, коротких поездок, документов**. Хорошо для **встреч, писем, контрактов, экзаменов**; **переговоры и обмен информацией** эффективны. Избегайте **крупных эмоциональных или финансовых решений**; **говорите и записывайте**.",
		"venus":   "Сейчас **час Венеры** — хорошо для **отношений, денег и красоты**. **Свидания, партнёрство, финансы, внешность, искусство** в потоке; **выражение симпатии, обсуждение условий** воспринимаются хорошо. Избегайте **конфронтации и больших рисков**; **гармония и приятные дела**.",
		"mars":    "Сейчас **час Марса** — энергия **действия и соревнования**. Хорошо для **спорта, рывков к дедлайнам, разрешения конфликтов, отстаивания прав**; **старт проектов, собеседования, состязания** уместны. Избегайте **бессмысленных споров и безрассудства**; **чёткость и решительность**.",
		"jupiter": "Сейчас **час Юпитера** — в пользу **расширения, возможностей и веры**. **Учёба, право, издательства, поездки, крупные планы** в потоке; **обращение за помощью, нетворкинг** получают отклик. Избегайте **перегрузки обязательствами и трат**; **плывите по течению, оставляйте запас**.",
		"saturn":  "Сейчас **час Сатурна** — хорошо для **долга, структуры и долгосрочных дел**. **Отчёты, обязательства, дисциплина, общение со старшими или правилами** в потоке; **завершение, разбор, планирование** эффективны. Избегайте **азарта и скоропалительных обещаний**; **практичность, шаг за шагом**.",
	}
	if s, ok := texts[string(planet)]; ok {
		return s
	}
	return "Текущий планетарный час подчёркивает сферы этой планеты; **планируйте соответствующие дела** для лучшего результата."
}

// ========== Void of Course (Moon) Interpretations (hourly) ==========

func getChineseVoidOfCourseInterpretation() string {
	return "月亮在换星座前未再与其它行星形成主要相位，即进入**月亮空亡**时段。传统上此时**不宜开启重要新事**（签约、表白、大额投资、手术择时等），因「事难成、易反复」。适合**收尾、休息、创意发想、日常琐事**；若必须行动，可当作**试水或排练**，留出后续调整空间。"
}

func getEnglishVoidOfCourseInterpretation() string {
	return "The Moon has left its last major aspect before changing sign — you're in a **void of course** period. Traditionally, **avoid starting major new things** (contracts, confessions, big investments, scheduling surgery) as outcomes may **stall or shift**. Good for **wrapping up, rest, brainstorming, routine tasks**; if you must act, treat it as a **trial run** and leave room to adjust."
}

func getRussianVoidOfCourseInterpretation() string {
	return "Луна вышла из последнего мажорного аспекта перед сменой знака — период **луны без курса**. Традиционно **не начинайте важных новых дел** (договоры, признания, крупные вложения, плановые операции): результаты могут **застопориться или измениться**. Подходит для **завершения дел, отдыха, мозгового штурма, рутины**; если нужно действовать, рассматривайте это как **пробный запуск** и оставьте место для корректировок."
}

// getRussianAspectInterpretation is in aspect_interpretations_ru.go

// ========== Void of Course (Moon) Interpretations (hourly) ==========

// ========== Lunar Phase Interpretations ==========

func (t *Translator) getLunarPhaseInterpretation(phase string) string {
	switch t.lang {
	case Chinese:
		return getChineseLunarPhaseInterpretation(phase)
	case Russian:
		return getRussianLunarPhaseInterpretation(phase)
	default:
		return getEnglishLunarPhaseInterpretation(phase)
	}
}

func getChineseLunarPhaseInterpretation(phase string) string {
	// Compatibility mapping for snake_case from astro
	if phase == "first_quarter" {
		phase = "firstQuarter"
	} else if phase == "last_quarter" {
		phase = "lastQuarter"
	}

	interpretations := map[string]string{
		"new": "新月如种子破土，这是**设定意图**的黄金时刻。宇宙能量正在**重新集结**，你内心深处的**新想法、新计划**会在这时萌芽。不要急于行动，先**静心思考**：未来一个月你最想实现什么？把愿望写下来，让它们在黑暗中**悄悄生长**。适合**冥想、许愿、规划**，避免**匆忙开始**。",
		"crescent": "新月后的**初生月牙**，能量开始**向上涌动**。就像幼苗破土而出，你需要**主动行动**来推动计划。这是**突破障碍**的时期，可能会遇到一些**阻力或质疑**，但正是这些挑战让你更坚定。不要退缩，**小步快跑**，每个微小的进展都在为未来铺路。适合**开始执行、克服困难、建立习惯**。",
		"firstQuarter": "**上弦月**带来**第一次考验**。你的计划可能遇到**现实的冲击**，需要**调整方向**或**做出选择**。这是**危机与机遇并存**的时刻，压力会促使你**重新审视目标**。不要害怕冲突，**主动沟通**，**果断决策**。那些经不起考验的想法会被淘汰，真正有价值的东西会**更加清晰**。适合**解决冲突、重新规划、做出决定**。",
		"gibbous": "**渐盈月**阶段，你的计划正在**逐步完善**。就像果实慢慢成熟，你需要**精雕细琢**，**查漏补缺**。这是**准备阶段**，不要急于求成，**耐心打磨细节**。可能会发现一些**需要改进的地方**，这是好事。**收集反馈**，**优化方案**，为即将到来的**收获期**做好准备。适合**完善细节、收集信息、调整策略**。",
		"full": "**满月**如**能量巅峰**，这是**收获与觉醒**的时刻。你之前种下的种子现在**开花结果**，一些事情会**达到高潮**或**真相大白**。情绪可能**更加敏感**，**直觉增强**。这是**庆祝成就**的时候，也是**释放旧模式**的机会。不要压抑情绪，**表达真实感受**，**放下不再服务你的东西**。适合**庆祝、表达、释放、做决定**。",
		"disseminating": "**渐亏月**开始，能量从**向外扩张**转向**向内整合**。这是**分享智慧**的时期，你积累的经验和知识可以**传递给他人**。不要吝啬，**教授、分享、传播**你的见解。同时开始**反思过程**，**总结经验教训**。适合**教学、分享、写作、传播理念**，避免**过度消耗**。",
		"lastQuarter": "**下弦月**带来**第二次考验**，这次是关于**放下与调整**。你需要**释放不再需要的东西**，**调整方向**，**告别旧模式**。可能会感到一些**失落或不适**，但这是**成长的必要过程**。不要执着于过去，**勇敢放手**，为新的周期**腾出空间**。适合**清理、释放、反思、调整**。",
		"balsamic": "**残月**如**最后的整合**，这是**休息与准备**的时期。能量**回归内在**，你需要**静心整合**这一周期的**所有经验**。不要急于开始新事物，**让一切沉淀**，**吸收养分**。这是**灵性连接**的时刻，**梦境和直觉**会带来重要信息。适合**休息、冥想、整合、准备**，避免**匆忙行动**。",
	}
	if text, ok := interpretations[phase]; ok {
		return text
	}
	return "月相变化带来**能量周期的转换**，注意观察**内在感受**和**外在机遇**的呼应。"
}

func getEnglishLunarPhaseInterpretation(phase string) string {
	// Compatibility mapping for snake_case from astro
	if phase == "first_quarter" {
		phase = "firstQuarter"
	} else if phase == "last_quarter" {
		phase = "lastQuarter"
	}

	interpretations := map[string]string{
		"new": "The **New Moon** is like a seed breaking ground — this is the **golden moment for setting intentions**. Cosmic energy is **regathering**, and **new ideas and plans** deep within you will sprout now. Don't rush into action; first **quietly reflect**: what do you most want to achieve in the coming month? Write down your wishes and let them **grow quietly in the dark**. Good for **meditation, wishing, planning**; avoid **hasty starts**.",
		"crescent": "The **Crescent Moon** after new moon, energy begins to **surge upward**. Like a seedling breaking through soil, you need **active steps** to move plans forward. This is a time of **breaking through obstacles**; you may encounter **resistance or doubt**, but these challenges make you more determined. Don't retreat; **take small, quick steps** — each tiny progress paves the way. Good for **starting execution, overcoming difficulties, building habits**.",
		"firstQuarter": "The **First Quarter Moon** brings the **first test**. Your plans may face **reality's impact**, requiring **course correction** or **making choices**. This is a moment of **crisis and opportunity**; pressure will push you to **reexamine goals**. Don't fear conflict; **communicate actively**, **decide decisively**. Ideas that can't stand the test will fall away; what truly matters becomes **clearer**. Good for **resolving conflicts, replanning, making decisions**.",
		"gibbous": "The **Gibbous Moon** phase — your plans are **gradually refining**. Like fruit slowly ripening, you need to **polish carefully**, **fill gaps**. This is the **preparation stage**; don't rush, **patiently refine details**. You may find **areas needing improvement** — that's good. **Gather feedback**, **optimize plans**, prepare for the coming **harvest**. Good for **perfecting details, gathering information, adjusting strategy**.",
		"full": "The **Full Moon** is like an **energy peak** — this is the moment of **harvest and awakening**. Seeds you planted earlier now **bloom and bear fruit**; some things reach **climax** or **truth emerges**. Emotions may be **more sensitive**, **intuition heightened**. This is time to **celebrate achievements** and **release old patterns**. Don't suppress emotions; **express true feelings**, **let go of what no longer serves**. Good for **celebration, expression, release, decision-making**.",
		"disseminating": "The **Disseminating Moon** begins; energy shifts from **outward expansion** to **inward integration**. This is the time to **share wisdom**; experience and knowledge you've gathered can **be passed to others**. Don't be stingy; **teach, share, spread** your insights. Also begin **reflecting on the process**, **summarizing lessons**. Good for **teaching, sharing, writing, spreading ideas**; avoid **overconsumption**.",
		"lastQuarter": "The **Last Quarter Moon** brings the **second test** — this time about **letting go and adjusting**. You need to **release what's no longer needed**, **adjust direction**, **say goodbye to old patterns**. You may feel some **loss or discomfort**, but this is **necessary for growth**. Don't cling to the past; **bravely let go**, **make space** for the new cycle. Good for **clearing, releasing, reflecting, adjusting**.",
		"balsamic": "The **Balsamic Moon** is like **final integration** — this is a time of **rest and preparation**. Energy **returns inward**; you need to **quietly integrate all experiences** from this cycle. Don't rush into new things; **let everything settle**, **absorb nourishment**. This is a moment of **spiritual connection**; **dreams and intuition** will bring important messages. Good for **rest, meditation, integration, preparation**; avoid **hasty action**.",
	}
	if text, ok := interpretations[phase]; ok {
		return text
	}
	return "Lunar phase changes bring **shifts in energy cycles**; notice how **inner feelings** and **outer opportunities** resonate."
}

func getRussianLunarPhaseInterpretation(phase string) string {
	// Compatibility mapping for snake_case from astro
	if phase == "first_quarter" {
		phase = "firstQuarter"
	} else if phase == "last_quarter" {
		phase = "lastQuarter"
	}

	interpretations := map[string]string{
		"new": "**Новолуние** как семя, пробивающееся сквозь землю — это **золотой момент для намерений**. Космическая энергия **перегруппировывается**, **новые идеи и планы** глубоко внутри вас прорастут сейчас. Не спешите действовать; сначала **тихо поразмышляйте**: чего вы больше всего хотите достичь в ближайший месяц? Запишите желания и позвольте им **тихо расти в темноте**. Хорошо для **медитации, желаний, планирования**; избегайте **поспешных стартов**.",
		"crescent": "**Молодая луна** после новолуния, энергия начинает **подниматься вверх**. Как росток, пробивающийся сквозь почву, вам нужны **активные шаги** для продвижения планов. Это время **прорыва через препятствия**; вы можете столкнуться с **сопротивлением или сомнениями**, но эти вызовы делают вас более решительными. Не отступайте; **делайте маленькие быстрые шаги** — каждый крошечный прогресс прокладывает путь. Хорошо для **начала исполнения, преодоления трудностей, построения привычек**.",
		"firstQuarter": "**Первая четверть** приносит **первое испытание**. Ваши планы могут столкнуться с **ударом реальности**, требуя **коррекции курса** или **принятия решений**. Это момент **кризиса и возможности**; давление заставит вас **пересмотреть цели**. Не бойтесь конфликта; **общайтесь активно**, **решайте решительно**. Идеи, которые не выдержат испытания, отпадут; то, что действительно важно, станет **яснее**. Хорошо для **разрешения конфликтов, перепланирования, принятия решений**.",
		"gibbous": "Фаза **растущей луны** — ваши планы **постепенно совершенствуются**. Как плод, медленно созревающий, вам нужно **тщательно полировать**, **заполнять пробелы**. Это **стадия подготовки**; не спешите, **терпеливо совершенствуйте детали**. Вы можете найти **области, требующие улучшения** — это хорошо. **Собирайте обратную связь**, **оптимизируйте планы**, готовьтесь к грядущему **урожаю**. Хорошо для **совершенствования деталей, сбора информации, корректировки стратегии**.",
		"full": "**Полнолуние** как **пик энергии** — это момент **урожая и пробуждения**. Семена, которые вы посадили ранее, теперь **цветут и плодоносят**; некоторые вещи достигают **кульминации** или **истина выходит наружу**. Эмоции могут быть **более чувствительными**, **интуиция обострена**. Это время **праздновать достижения** и **освобождаться от старых паттернов**. Не подавляйте эмоции; **выражайте истинные чувства**, **отпускайте то, что больше не служит**. Хорошо для **празднования, выражения, освобождения, принятия решений**.",
		"disseminating": "**Убывающая луна** начинается; энергия переходит от **внешнего расширения** к **внутренней интеграции**. Это время **делиться мудростью**; опыт и знания, которые вы накопили, могут **быть переданы другим**. Не будьте скупы; **учите, делитесь, распространяйте** свои идеи. Также начинайте **размышлять о процессе**, **подводить итоги уроков**. Хорошо для **обучения, обмена, письма, распространения идей**; избегайте **чрезмерного потребления**.",
		"lastQuarter": "**Последняя четверть** приносит **второе испытание** — на этот раз о **отпускании и корректировке**. Вам нужно **освободиться от ненужного**, **скорректировать направление**, **попрощаться со старыми паттернами**. Вы можете почувствовать некоторую **потерю или дискомфорт**, но это **необходимо для роста**. Не цепляйтесь за прошлое; **смело отпускайте**, **освобождайте место** для нового цикла. Хорошо для **очищения, освобождения, размышления, корректировки**.",
		"balsamic": "**Бальзамическая луна** как **финальная интеграция** — это время **отдыха и подготовки**. Энергия **возвращается внутрь**; вам нужно **тихо интегрировать весь опыт** этого цикла. Не спешите начинать новое; **позвольте всему осесть**, **впитывайте питание**. Это момент **духовной связи**; **сны и интуиция** принесут важные сообщения. Хорошо для **отдыха, медитации, интеграции, подготовки**; избегайте **поспешных действий**.",
	}
	if text, ok := interpretations[phase]; ok {
		return text
	}
	return "Изменения лунных фаз приносят **сдвиги в энергетических циклах**; замечайте, как **внутренние чувства** и **внешние возможности** резонируют."
}

// ========== Sign Change Interpretations ==========

func (t *Translator) getSignChangeInterpretation(planet models.PlanetID, newSign string) string {
	switch t.lang {
	case Chinese:
		return getChineseSignChangeInterpretation(planet, newSign)
	case Russian:
		return getRussianSignChangeInterpretation(planet, newSign)
	default:
		return getEnglishSignChangeInterpretation(planet, newSign)
	}
}

func getChineseSignChangeInterpretation(planet models.PlanetID, newSign string) string {
	key := string(planet) + "_" + newSign
	interpretations := map[string]string{
		// Sun sign changes - 每个星座都有独特的表达
		"sun_aries": "太阳进入**白羊座**，**行动力**和**开创精神**被点燃。你会感到**能量充沛**，想要**立即行动**，**不再等待**。这是**启动新项目**的好时机，但要注意**避免冲动**，**三思而后行**。适合**竞争、运动、领导**，避免**过于急躁或独断**。",
		"sun_taurus": "太阳进入**金牛座**，**稳定**和**享受**成为主题。你会**放慢节奏**，**品味生活**，**重视物质安全**。这是**建立基础**、**积累资源**的好时期，但要注意**不要过于固执**。适合**理财、艺术、美食**，避免**过度消费或懒惰**。",
		"sun_gemini": "太阳进入**双子座**，**沟通**和**学习**变得活跃。你会**思维敏捷**，**好奇心强**，**想要交流**。这是**学习新知识**、**建立人脉**的好时机，但要注意**避免浅尝辄止**。适合**写作、教学、短途旅行**，避免**过于分散注意力**。",
		"sun_cancer": "太阳进入**巨蟹座**，**情感**和**家庭**成为焦点。你会**更加敏感**，**重视安全感**，**想要照顾他人**。这是**改善家庭关系**、**处理房产**的好时期，但要注意**不要过度情绪化**。适合**家庭活动、烹饪、情感交流**，避免**过度依赖或封闭**。",
		"sun_leo": "太阳进入**狮子座**，**创造力**和**表现力**达到高峰。你会**自信满满**，**想要展现自己**，**享受关注**。这是**艺术创作**、**娱乐活动**的好时机，但要注意**不要过于自我中心**。适合**表演、创作、领导**，避免**过度炫耀或自负**。",
		"sun_virgo": "太阳进入**处女座**，**细节**和**服务**被强调。你会**更加谨慎**，**追求完美**，**想要改善**。这是**优化工作流程**、**健康管理**的好时期，但要注意**不要过于挑剔**。适合**整理、分析、服务他人**，避免**过度批评或焦虑**。",
		"sun_libra": "太阳进入**天秤座**，**平衡**和**关系**成为主题。你会**重视和谐**，**追求美感**，**想要合作**。这是**改善关系**、**艺术欣赏**的好时机，但要注意**不要过度妥协**。适合**社交、艺术、谈判**，避免**优柔寡断或依赖他人**。",
		"sun_scorpio": "太阳进入**天蝎座**，**深度**和**转化**被激活。你会**更加专注**，**深入探索**，**想要改变**。这是**处理深层问题**、**资源管理**的好时期，但要注意**不要过于极端**。适合**研究、投资、心理工作**，避免**过度控制或报复**。",
		"sun_sagittarius": "太阳进入**射手座**，**自由**和**探索**成为焦点。你会**视野开阔**，**乐观向上**，**想要冒险**。这是**学习、旅行、追求真理**的好时机，但要注意**不要过于理想化**。适合**教育、旅行、哲学**，避免**过度承诺或不负责任**。",
		"sun_capricorn": "太阳进入**摩羯座**，**责任**和**成就**被强调。你会**更加务实**，**重视结构**，**想要建立**。这是**职业发展**、**长期规划**的好时期，但要注意**不要过于严肃**。适合**工作、规划、建立权威**，避免**过度工作或冷漠**。",
		"sun_aquarius": "太阳进入**水瓶座**，**创新**和**独立**成为主题。你会**思维独特**，**重视自由**，**想要改变**。这是**科技、社群、理想**的好时机，但要注意**不要过于疏离**。适合**创新、社交、追求理想**，避免**过度叛逆或冷漠**。",
		"sun_pisces": "太阳进入**双鱼座**，**直觉**和**慈悲**被激活。你会**更加敏感**，**富有想象力**，**想要连接**。这是**艺术、灵性、疗愈**的好时期，但要注意**不要过于逃避**。适合**创作、冥想、帮助他人**，避免**过度幻想或牺牲**。",
		
		// Moon sign changes - 更快速的情感变化
		"moon_aries": "月亮进入**白羊座**，**情绪**变得**直接而冲动**。你会**想要立即行动**，**表达愤怒或热情**，**不再压抑**。这是**释放情绪**的好时机，但要注意**控制脾气**。适合**运动、竞争、快速决策**。",
		"moon_taurus": "月亮进入**金牛座**，**情绪**变得**稳定而舒适**。你会**想要享受**，**寻求安全感**，**放慢节奏**。这是**放松、美食、艺术**的好时期。适合**休息、购物、感官享受**。",
		"moon_gemini": "月亮进入**双子座**，**情绪**变得**活跃而多变**。你会**想要交流**，**学习新东西**，**思维跳跃**。这是**社交、学习、短途出行**的好时机。适合**聊天、阅读、短途旅行**。",
		"moon_cancer": "月亮进入**巨蟹座**，**情绪**变得**敏感而保护**。你会**想要回家**，**照顾他人**，**情感丰富**。这是**家庭时间、情感交流**的好时期。适合**家庭活动、烹饪、情感表达**。",
		"moon_leo": "月亮进入**狮子座**，**情绪**变得**自信而热情**。你会**想要展现**，**享受关注**，**创造快乐**。这是**娱乐、创作、表演**的好时机。适合**艺术、娱乐、领导**。",
		"moon_virgo": "月亮进入**处女座**，**情绪**变得**谨慎而服务**。你会**想要改善**，**关注细节**，**帮助他人**。这是**整理、分析、健康管理**的好时期。适合**工作、整理、服务**。",
		"moon_libra": "月亮进入**天秤座**，**情绪**变得**和谐而平衡**。你会**想要合作**，**追求美感**，**避免冲突**。这是**社交、艺术、关系**的好时机。适合**约会、艺术、谈判**。",
		"moon_scorpio": "月亮进入**天蝎座**，**情绪**变得**深刻而强烈**。你会**想要深入**，**探索真相**，**情感强烈**。这是**深度交流、研究、转化**的好时期。适合**深度对话、研究、心理工作**。",
		"moon_sagittarius": "月亮进入**射手座**，**情绪**变得**乐观而自由**。你会**想要冒险**，**探索新事物**，**视野开阔**。这是**学习、旅行、哲学**的好时机。适合**教育、旅行、追求理想**。",
		"moon_capricorn": "月亮进入**摩羯座**，**情绪**变得**严肃而克制**。你会**想要控制**，**承担责任**，**建立结构**。这是**工作、规划、建立权威**的好时期。适合**工作、规划、建立目标**。",
		"moon_aquarius": "月亮进入**水瓶座**，**情绪**变得**独立而理性**。你会**想要自由**，**追求理想**，**保持距离**。这是**创新、社交、理想**的好时机。适合**创新、团体活动、追求理想**。",
		"moon_pisces": "月亮进入**双鱼座**，**情绪**变得**敏感而梦幻**。你会**想要连接**，**富有同情心**，**直觉增强**。这是**艺术、灵性、疗愈**的好时期。适合**创作、冥想、帮助他人**。",
	}
	if text, ok := interpretations[key]; ok {
		return text
	}
	// Fallback: 根据行星和星座生成通用解释
	planetName := map[string]string{
		"sun": "太阳", "moon": "月亮", "mercury": "水星", "venus": "金星", "mars": "火星",
		"jupiter": "木星", "saturn": "土星", "uranus": "天王星", "neptune": "海王星", "pluto": "冥王星",
	}[string(planet)]
	return planetName + "进入**" + newSign + "座**，带来该星座特质的能量。注意观察**相关生活领域**的变化，**顺应能量**，**善用优势**。"
}

func getEnglishSignChangeInterpretation(planet models.PlanetID, newSign string) string {
	key := string(planet) + "_" + newSign
	interpretations := map[string]string{
		"sun_aries": "The Sun enters **Aries**, igniting **action** and **pioneering spirit**. You'll feel **energetic**, wanting to **act immediately**, **no more waiting**. Good time to **launch new projects**, but watch **avoiding impulsiveness**; **think before acting**. Suitable for **competition, sports, leadership**; avoid **being too hasty or dictatorial**.",
		"sun_taurus": "The Sun enters **Taurus**; **stability** and **enjoyment** become themes. You'll **slow down**, **savor life**, **value material security**. Good period for **building foundations**, **accumulating resources**, but watch **not being too stubborn**. Suitable for **finance, arts, food**; avoid **overspending or laziness**.",
		"sun_gemini": "The Sun enters **Gemini**; **communication** and **learning** become active. You'll be **quick-witted**, **curious**, **wanting to connect**. Good time for **learning**, **networking**, but watch **avoiding superficiality**. Suitable for **writing, teaching, short trips**; avoid **scattering attention**.",
		"sun_cancer": "The Sun enters **Cancer**; **emotions** and **family** become focus. You'll be **more sensitive**, **value security**, **want to care for others**. Good period for **improving family relations**, **property matters**, but watch **not being overly emotional**. Suitable for **family activities, cooking, emotional exchange**; avoid **overdependence or withdrawal**.",
		"sun_leo": "The Sun enters **Leo**; **creativity** and **expression** peak. You'll be **confident**, **wanting to shine**, **enjoy attention**. Good time for **artistic creation**, **entertainment**, but watch **not being self-centered**. Suitable for **performance, creation, leadership**; avoid **excessive showiness or arrogance**.",
		"sun_virgo": "The Sun enters **Virgo**; **detail** and **service** are emphasized. You'll be **more cautious**, **pursuing perfection**, **wanting to improve**. Good period for **optimizing work**, **health management**, but watch **not being too critical**. Suitable for **organizing, analyzing, serving others**; avoid **excessive criticism or anxiety**.",
		"sun_libra": "The Sun enters **Libra**; **balance** and **relationships** become themes. You'll **value harmony**, **pursue beauty**, **want to cooperate**. Good time for **improving relationships**, **art appreciation**, but watch **not overcompromising**. Suitable for **socializing, arts, negotiation**; avoid **indecisiveness or dependence**.",
		"sun_scorpio": "The Sun enters **Scorpio**; **depth** and **transformation** activate. You'll be **more focused**, **exploring deeply**, **wanting change**. Good period for **dealing with deep issues**, **resource management**, but watch **not being extreme**. Suitable for **research, investment, psychological work**; avoid **excessive control or revenge**.",
		"sun_sagittarius": "The Sun enters **Sagittarius**; **freedom** and **exploration** become focus. You'll have **broad vision**, **optimism**, **wanting adventure**. Good time for **learning, travel, pursuing truth**, but watch **not being too idealistic**. Suitable for **education, travel, philosophy**; avoid **overcommitting or irresponsibility**.",
		"sun_capricorn": "The Sun enters **Capricorn**; **responsibility** and **achievement** are emphasized. You'll be **more practical**, **value structure**, **wanting to build**. Good period for **career development**, **long-term planning**, but watch **not being too serious**. Suitable for **work, planning, establishing authority**; avoid **overwork or coldness**.",
		"sun_aquarius": "The Sun enters **Aquarius**; **innovation** and **independence** become themes. You'll have **unique thinking**, **value freedom**, **wanting change**. Good time for **technology, community, ideals**, but watch **not being too detached**. Suitable for **innovation, socializing, pursuing ideals**; avoid **excessive rebellion or coldness**.",
		"sun_pisces": "The Sun enters **Pisces**; **intuition** and **compassion** activate. You'll be **more sensitive**, **imaginative**, **wanting connection**. Good period for **arts, spirituality, healing**, but watch **not escaping**. Suitable for **creation, meditation, helping others**; avoid **excessive fantasy or sacrifice**.",
	}
	if text, ok := interpretations[key]; ok {
		return text
	}
	planetName := map[string]string{
		"sun": "Sun", "moon": "Moon", "mercury": "Mercury", "venus": "Venus", "mars": "Mars",
		"jupiter": "Jupiter", "saturn": "Saturn", "uranus": "Uranus", "neptune": "Neptune", "pluto": "Pluto",
	}[string(planet)]
	return planetName + " enters **" + newSign + "**, bringing that sign's energy. Notice changes in **related life areas**; **flow with energy**, **use advantages**."
}

func getRussianSignChangeInterpretation(planet models.PlanetID, newSign string) string {
	key := string(planet) + "_" + newSign
	interpretations := map[string]string{
		"sun_aries":       "Солнце входит в **Овен**, зажигая **действие** и **пионерский дух**. Вы почувствуете **энергию**, захотите **действовать немедленно**, **больше не ждать**. Хорошее время для **запуска новых проектов**, но следите за **избеганием импульсивности**; **думайте перед действием**. Подходит для **соревнований, спорта, лидерства**; избегайте **слишком поспешного или диктаторского**.",
		"sun_taurus":      "Солнце входит в **Телец**; **стабильность** и **удовольствие** становятся темами. Вы **замедлитесь**, **будете наслаждаться жизнью**, **ценить материальную безопасность**. Хороший период для **построения основ**, **накопления ресурсов**, но следите за **не быть слишком упрямым**. Подходит для **финансов, искусства, еды**; избегайте **чрезмерных трат или лени**.",
		"sun_gemini":      "Солнце входит в **Близнецы**; **общение** и **обучение** становятся активными. Вы будете **быстро соображать**, **любопытны**, **хотеть общаться**. Хорошее время для **обучения**, **сетей**, но следите за **избеганием поверхностности**. Подходит для **письма, преподавания, коротких поездок**; избегайте **рассеивания внимания**.",
		"sun_cancer":      "Солнце входит в **Рак**; **эмоции** и **семья** становятся центром внимания. Вы станете **чувствительнее**, будете **ценить безопасность**, захотите **заботиться о других**. Хороший период для **улучшения семейных отношений**, **дел с недвижимостью**, но следите за тем, чтобы **не быть слишком эмоциональным**. Подходит для **семейных мероприятий, кулинарии, эмоционального обмена**.",
		"sun_leo":         "Солнце входит в **Лев**; **творчество** и **самовыражение** достигают пика. Вы будете **уверены в себе**, захотите **сиять**, **наслаждаться вниманием**. Хорошее время для **художественного творчества**, **развлечений**, но следите за тем, чтобы **не быть эгоцентричным**. Подходит для **выступлений, творчества, лидерства**.",
		"sun_virgo":       "Солнце входит в **Дева**; подчеркиваются **детали** и **служение**. Вы станете **осторожнее**, будете **стремиться к совершенству**, захотите **улучшений**. Хороший период для **оптимизации работы**, **управления здоровьем**, но следите за тем, чтобы **не быть слишком критичным**. Подходит для **организации, анализа, помощи другим**.",
		"sun_libra":       "Солнце входит в **Весы**; **баланс** и **отношения** становятся темами. Вы будете **ценить гармонию**, **стремиться к красоте**, захотите **сотрудничать**. Хорошее время для **улучшения отношений**, **оценки искусства**, но следите за тем, чтобы **не идти на чрезмерные компромиссы**. Подходит для **общения, искусства, переговоров**.",
		"sun_scorpio":     "Солнце входит в **Скорпион**; активируются **глубина** и **трансформация**. Вы станете **сосредоточеннее**, будете **глубоко исследовать**, захотите **перемен**. Хороший период для **решения глубоких проблем**, **управления ресурсами**, но следите за тем, чтобы **не впадать в крайности**. Подходит для **исследований, инвестиций, психологической работы**.",
		"sun_sagittarius": "Солнце входит в **Стрелец**; **свобода** и **исследование** становятся центром внимания. У вас будет **широкое видение**, **оптимизм**, желание **приключений**. Хорошее время для **обучения, путешествий, поиска истины**, но следите за тем, чтобы **не быть слишком идеалистичным**. Подходит для **образования, путешествий, философии**.",
		"sun_capricorn":   "Солнце входит в **Козерог**; подчеркиваются **ответственность** и **достижения**. Вы станете **практичнее**, будете **ценить структуру**, захотите **строить**. Хороший период для **развития карьеры**, **долгосрочного планирования**, но следите за тем, чтобы **не быть слишком серьезным**. Подходит для **работы, планирования, укрепления авторитета**.",
		"sun_aquarius":    "Солнце входит в **Водолей**; **инновации** и **независимость** становятся темами. У вас будет **уникальное мышление**, вы будете **ценить свободу**, захотите **перемен**. Хорошее время для **технологий, сообществ, идеалов**, но следите за тем, чтобы **не быть слишком отстраненным**. Подходит для **инноваций, общения, погони за идеалами**.",
		"sun_pisces":      "Солнце входит в **Рыбы**; активируются **интуиция** и **сострадание**. Вы станете **чувствительнее**, **воображение** усилится, возникнет желание **связи**. Хороший период для **искусства, духовности, исцеления**, но следите за тем, чтобы **не убегать от реальности**. Подходит для **творчества, медитации, помощи другим**.",
		"moon_aries":      "Луна входит в **Овен**; **эмоции** становятся **прямыми и импульсивными**. Вы захотите **действовать немедленно**, **выражать гнев или энтузиазм**, **больше не подавлять**. Хорошее время для **выплеска эмоций**, но следите за **темпераментом**. Подходит для **спорта, соревнований, быстрых решений**.",
		"moon_taurus":     "Луна входит в **Телец**; **эмоции** становятся **стабильными и комфортными**. Вы захотите **наслаждаться**, **искать безопасность**, **замедлить темп**. Хороший период для **отдыха, вкусной еды, искусства**. Подходит для **отдыха, покупок, чувственных удовольствий**.",
		"moon_gemini":     "Луна входит в **Близнецы**; **эмоции** становятся **активными и переменчивыми**. Вы захотите **общаться**, **узнавать новое**, **мысли будут перескакивать**. Хорошее время для **общения, обучения, коротких поездок**. Подходит для **разговоров, чтения, путешествий**.",
		"moon_cancer":     "Луна входит в **Рак**; **эмоции** становятся **чувствительными и защитными**. Вы захотите **домой**, **заботиться о других**, **чувства будут богатыми**. Хороший период для **семейного времени, эмоционального обмена**. Подходит для **семейных дел, кулинарии, выражения чувств**.",
		"moon_leo":        "Луна входит в **Лев**; **эмоции** становятся **уверенными и страстными**. Вы захотите **показать себя**, **наслаждаться вниманием**, **создавать радость**. Хорошее время для **развлечений, творчества, выступлений**. Подходит для **искусства, развлечений, лидерства**.",
		"moon_virgo":      "Луна входит в **Дева**; **эмоции** становятся **осторожными и служебными**. Вы захотите **улучшений**, будете **внимательны к деталям**, захотите **помогать другим**. Хороший период для **организации, анализа, управления здоровьем**. Подходит для **работы, уборки, служения**.",
		"moon_libra":      "Луна входит в **Весы**; **эмоции** становятся **гармоничными и сбалансированными**. Вы захотите **сотрудничать**, **стремиться к красоте**, **избегать конфликтов**. Хорошее время для **общения, искусства, отношений**. Подходит для **свиданий, искусства, переговоров**.",
		"moon_scorpio":    "Луна входит в **Скорпион**; **эмоции** становятся **глубокими и интенсивными**. Вы захотите **глубины**, **поиска истины**, **чувства будут сильными**. Хороший период для **глубокого общения, исследований, трансформации**. Подходит для **глубоких разговоров, исследований, психологической работы**.",
		"moon_sagittarius": "Луна входит в **Стрелец**; **эмоции** становятся **оптимистичными и свободными**. Вы захотите **приключений**, **исследования нового**, **видение расширится**. Хорошее время для **обучения, путешествий, философии**. Подходит для **образования, путешествий, погони за идеалами**.",
		"moon_capricorn":   "Луна входит в **Козерог**; **эмоции** становятся **серьезными и сдержанными**. Вы захотите **контроля**, **принятия ответственности**, **создания структуры**. Хороший период для **работы, планирования, укрепления авторитета**. Подходит для **работы, планирования, постановки целей**.",
		"moon_aquarius":    "Луна входит в **Водолей**; **эмоции** становятся **независимыми и рациональными**. Вы захотите **свободы**, **погони за идеалами**, **сохранения дистанции**. Хорошее время для **инноваций, групповой деятельности, идеалов**. Подходит для **инноваций, социальных мероприятий, погони за идеалами**.",
		"moon_pisces":      "Луна входит в **Рыбы**; **эмоции** становятся **чувствительными и мечтательными**. Вы захотите **связи**, будете **сострадательны**, **интуиция усилится**. Хороший период для **искусства, духовности, исцеления**. Подходит для **творчества, медитации, помощи другим**.",
	}
	if text, ok := interpretations[key]; ok {
		return text
	}
	planetName := map[string]string{
		"sun": "Солнце", "moon": "Луна", "mercury": "Меркурий", "venus": "Венера", "mars": "Марс",
		"jupiter": "Юпитер", "saturn": "Сатурн", "uranus": "Уран", "neptune": "Нептун", "pluto": "Плутон",
	}[string(planet)]
	return planetName + " входит в **" + newSign + "**, принося энергию этого знака. Замечайте изменения в **связанных сферах жизни**; **плывите с энергией**, **используйте преимущества**."
}

// ========== Dignity Interpretations ==========

func (t *Translator) getDignityInterpretation(planet models.PlanetID, dignityType string) string {
	switch t.lang {
	case Chinese:
		return getChineseDignityInterpretation(planet, dignityType)
	case Russian:
		return getRussianDignityInterpretation(planet, dignityType)
	default:
		return getEnglishDignityInterpretation(planet, dignityType)
	}
}

func getChineseDignityInterpretation(planet models.PlanetID, dignityType string) string {
	key := string(planet) + "_" + dignityType
	interpretations := map[string]string{
		// Domicile (入庙) - 行星在自己守护的星座
		"sun_domicile": "太阳在**狮子座**（入庙），**自我表达**和**领导力**达到**自然巅峰**。你会**自信满满**，**光芒四射**，**吸引他人注意**。这是**展现才华**、**承担领导角色**的绝佳时期。能量**稳定而强大**，适合**创作、表演、领导项目**。",
		"moon_domicile": "月亮在**巨蟹座**（入庙），**情感**和**直觉**处于**最舒适的状态**。你会**情感丰富**，**保护性强**，**家庭意识强烈**。这是**处理家庭事务**、**情感疗愈**的好时期。能量**滋养而安全**，适合**照顾他人、家居装饰、情感表达**。",
		"mercury_domicile": "水星在**双子座**或**处女座**（入庙），**思维**和**沟通**达到**最佳状态**。你会**思维敏捷**，**表达清晰**，**学习能力强**。这是**学习、写作、谈判**的绝佳时期。能量**灵活而精确**，适合**教学、分析、信息处理**。",
		"venus_domicile": "金星在**金牛座**或**天秤座**（入庙），**爱与美**处于**最和谐的状态**。你会**魅力四射**，**审美敏锐**，**关系和谐**。这是**深化关系**、**艺术创作**的好时期。能量**优雅而平衡**，适合**约会、艺术、享受生活**。",
		"mars_domicile": "火星在**白羊座**或**天蝎座**（入庙），**行动力**和**意志力**达到**最强状态**。你会**精力充沛**，**目标明确**，**执行力强**。这是**启动项目**、**克服障碍**的绝佳时期。能量**直接而有力**，适合**竞争、运动、追求目标**。",
		"jupiter_domicile": "木星在**射手座**或**双鱼座**（入庙），**扩展**和**智慧**处于**最理想的状态**。你会**视野开阔**，**乐观向上**，**机会增多**。这是**学习、旅行、追求理想**的好时期。能量**慷慨而智慧**，适合**教育、哲学、探索**。",
		"saturn_domicile": "土星在**摩羯座**或**水瓶座**（入庙），**责任**和**结构**达到**最稳定的状态**。你会**更加务实**，**纪律性强**，**长期规划清晰**。这是**建立事业**、**承担责任**的绝佳时期。能量**稳定而持久**，适合**工作、规划、建立权威**。",
		
		// Exaltation (旺相) - 行星在旺相的星座
		"sun_exaltation": "太阳在**白羊座**（旺相），**开创精神**和**行动力**被**极大增强**。你会**充满活力**，**敢于冒险**，**领导力突出**。这是**启动新项目**、**展现勇气**的好时期。能量**积极而冲动**，适合**竞争、运动、领导**。",
		"moon_exaltation": "月亮在**金牛座**（旺相），**情感**和**安全感**处于**最舒适的状态**。你会**情绪稳定**，**享受生活**，**物质满足**。这是**建立安全感**、**享受生活**的好时期。能量**滋养而稳定**，适合**休息、美食、艺术**。",
		"mercury_exaltation": "水星在**处女座**（旺相），**分析能力**和**精确度**达到**最高水平**。你会**思维严谨**，**注重细节**，**效率极高**。这是**优化工作**、**健康管理**的绝佳时期。能量**精确而高效**，适合**分析、整理、服务**。",
		"venus_exaltation": "金星在**双鱼座**（旺相），**爱与慈悲**处于**最理想的状态**。你会**富有同情心**，**艺术敏感**，**情感丰富**。这是**艺术创作**、**情感连接**的好时期。能量**温柔而梦幻**，适合**创作、冥想、帮助他人**。",
		"mars_exaltation": "火星在**摩羯座**（旺相），**行动力**和**持久力**达到**最强状态**。你会**目标明确**，**执行力强**，**坚持不懈**。这是**建立事业**、**克服困难**的绝佳时期。能量**坚定而持久**，适合**工作、竞争、追求目标**。",
		"jupiter_exaltation": "木星在**巨蟹座**（旺相），**扩展**和**滋养**处于**最理想的状态**。你会**家庭幸福**，**情感丰富**，**机会增多**。这是**改善家庭**、**情感成长**的好时期。能量**滋养而扩展**，适合**家庭活动、情感交流、建立安全感**。",
		"saturn_exaltation": "土星在**天秤座**（旺相），**责任**和**平衡**达到**最和谐的状态**。你会**关系成熟**，**公正公平**，**结构清晰**。这是**改善关系**、**建立合作**的绝佳时期。能量**平衡而稳定**，适合**合作、谈判、建立结构**。",
		
		// Detriment (落陷) - 行星在相反守护的星座
		"sun_detriment": "太阳在**水瓶座**（落陷），**自我表达**可能**受到限制**。你会**更加独立**，**疏离感强**，**难以展现**。这是**重新定义自我**、**追求理想**的时期，但要注意**不要过于疏离**。能量**独特而疏离**，需要**平衡个人与集体**。",
		"moon_detriment": "月亮在**摩羯座**（落陷），**情感**可能**受到压抑**。你会**情感克制**，**严肃认真**，**难以表达**。这是**承担责任**、**建立结构**的时期，但要注意**不要过度压抑情感**。能量**严肃而克制**，需要**允许情感流动**。",
		"mercury_detriment": "水星在**射手座**或**双鱼座**（落陷），**思维**可能**不够精确**。你会**思维跳跃**，**理想化**，**容易分散**。这是**追求理想**、**直觉思考**的时期，但要注意**不要过于理想化**。能量**理想而分散**，需要**保持专注**。",
		"venus_detriment": "金星在**白羊座**或**天蝎座**（落陷），**爱与美**可能**表达困难**。你会**关系直接**，**情感强烈**，**难以平衡**。这是**追求激情**、**深度连接**的时期，但要注意**不要过于极端**。能量**强烈而直接**，需要**平衡与和谐**。",
		"mars_detriment": "火星在**金牛座**或**天秤座**（落陷），**行动力**可能**受到阻碍**。你会**行动缓慢**，**犹豫不决**，**难以启动**。这是**稳定基础**、**寻求平衡**的时期，但要注意**不要过于被动**。能量**稳定而缓慢**，需要**主动推动**。",
		"jupiter_detriment": "木星在**双子座**或**处女座**（落陷），**扩展**可能**受到限制**。你会**视野狭窄**，**过于谨慎**，**机会减少**。这是**学习细节**、**建立基础**的时期，但要注意**不要过于局限**。能量**谨慎而局限**，需要**扩大视野**。",
		"saturn_detriment": "土星在**巨蟹座**或**狮子座**（落陷），**责任**可能**表达困难**。你会**情感负担重**，**难以建立结构**，**权威感弱**。这是**处理情感**、**建立自信**的时期，但要注意**不要过度负担**。能量**情感而负担**，需要**建立边界**。",
		
		// Fall (失势) - 行星在相反旺相的星座
		"sun_fall": "太阳在**天秤座**（失势），**自我表达**可能**需要平衡**。你会**重视关系**，**难以独立**，**需要合作**。这是**改善关系**、**寻求平衡**的时期，但要注意**不要失去自我**。能量**平衡而依赖**，需要**保持独立**。",
		"moon_fall": "月亮在**天蝎座**（失势），**情感**可能**过于强烈**。你会**情感深刻**，**控制欲强**，**难以释放**。这是**深度转化**、**处理阴影**的时期，但要注意**不要过度控制**。能量**深刻而强烈**，需要**释放与信任**。",
		"mercury_fall": "水星在**双鱼座**（失势），**思维**可能**不够清晰**。你会**直觉增强**，**容易混淆**，**难以精确**。这是**直觉思考**、**艺术创作**的时期，但要注意**不要过于模糊**。能量**直觉而模糊**，需要**保持清晰**。",
		"venus_fall": "金星在**处女座**（失势），**爱与美**可能**表达困难**。你会**过于挑剔**，**难以享受**，**关系紧张**。这是**改善关系**、**追求完美**的时期，但要注意**不要过于批评**。能量**挑剔而紧张**，需要**接受不完美**。",
		"mars_fall": "火星在**巨蟹座**（失势），**行动力**可能**受到情绪影响**。你会**情绪驱动**，**行动被动**，**难以直接**。这是**处理情感**、**建立安全感**的时期，但要注意**不要过度情绪化**。能量**情绪而被动**，需要**理性行动**。",
		"jupiter_fall": "木星在**摩羯座**（失势），**扩展**可能**受到限制**。你会**过于谨慎**，**机会减少**，**视野狭窄**。这是**建立基础**、**承担责任**的时期，但要注意**不要过于局限**。能量**谨慎而局限**，需要**扩大视野**。",
		"saturn_fall": "土星在**白羊座**（失势），**责任**可能**表达困难**。你会**冲动行动**，**难以建立结构**，**缺乏耐心**。这是**快速行动**、**展现勇气**的时期，但要注意**不要过于冲动**。能量**冲动而缺乏结构**，需要**建立纪律**。",
	}
	if text, ok := interpretations[key]; ok {
		return text
	}
	planetName := map[string]string{
		"sun": "太阳", "moon": "月亮", "mercury": "水星", "venus": "金星", "mars": "火星",
		"jupiter": "木星", "saturn": "土星", "uranus": "天王星", "neptune": "海王星", "pluto": "冥王星",
	}[string(planet)]
	dignityName := map[string]string{
		"domicile": "入庙", "exaltation": "旺相", "detriment": "落陷", "fall": "失势",
	}[dignityType]
	return planetName + "处于**" + dignityName + "**状态，能量表达**" + (map[string]string{"domicile": "自然顺畅", "exaltation": "增强提升", "detriment": "受到限制", "fall": "需要调整"}[dignityType]) + "**。注意**善用优势**，**弥补不足**。"
}

func getEnglishDignityInterpretation(planet models.PlanetID, dignityType string) string {
	key := string(planet) + "_" + dignityType
	interpretations := map[string]string{
		"sun_domicile": "The Sun in **Leo** (domicile) — **self-expression** and **leadership** reach **natural peak**. You'll be **confident**, **radiant**, **attracting attention**. Excellent time for **showing talent**, **taking leadership roles**. Energy is **stable and powerful**; suitable for **creation, performance, leading projects**.",
		"moon_domicile": "The Moon in **Cancer** (domicile) — **emotions** and **intuition** are in **most comfortable state**. You'll be **emotionally rich**, **protective**, **strong family awareness**. Good period for **handling family matters**, **emotional healing**. Energy is **nurturing and safe**; suitable for **caring for others, home decoration, emotional expression**.",
		"mercury_domicile": "Mercury in **Gemini** or **Virgo** (domicile) — **thinking** and **communication** reach **best state**. You'll be **quick-witted**, **clear expression**, **strong learning ability**. Excellent time for **learning, writing, negotiation**. Energy is **flexible and precise**; suitable for **teaching, analysis, information processing**.",
		"venus_domicile": "Venus in **Taurus** or **Libra** (domicile) — **love and beauty** are in **most harmonious state**. You'll be **charming**, **aesthetically sharp**, **harmonious relationships**. Good period for **deepening relationships**, **artistic creation**. Energy is **elegant and balanced**; suitable for **dating, arts, enjoying life**.",
		"mars_domicile": "Mars in **Aries** or **Scorpio** (domicile) — **action** and **willpower** reach **strongest state**. You'll be **energetic**, **clear goals**, **strong execution**. Excellent time for **launching projects**, **overcoming obstacles**. Energy is **direct and powerful**; suitable for **competition, sports, pursuing goals**.",
		"jupiter_domicile": "Jupiter in **Sagittarius** or **Pisces** (domicile) — **expansion** and **wisdom** are in **most ideal state**. You'll have **broad vision**, **optimism**, **increased opportunities**. Good period for **learning, travel, pursuing ideals**. Energy is **generous and wise**; suitable for **education, philosophy, exploration**.",
		"saturn_domicile": "Saturn in **Capricorn** or **Aquarius** (domicile) — **responsibility** and **structure** reach **most stable state**. You'll be **more practical**, **disciplined**, **clear long-term planning**. Excellent time for **building career**, **taking responsibility**. Energy is **stable and lasting**; suitable for **work, planning, establishing authority**.",
		"sun_exaltation": "The Sun in **Aries** (exaltation) — **pioneering spirit** and **action** are **greatly enhanced**. You'll be **full of vitality**, **daring to adventure**, **outstanding leadership**. Good time for **launching new projects**, **showing courage**. Energy is **active and impulsive**; suitable for **competition, sports, leadership**.",
		"moon_exaltation": "The Moon in **Taurus** (exaltation) — **emotions** and **security** are in **most comfortable state**. You'll be **emotionally stable**, **enjoying life**, **material satisfaction**. Good period for **building security**, **enjoying life**. Energy is **nurturing and stable**; suitable for **rest, food, arts**.",
		"mercury_exaltation": "Mercury in **Virgo** (exaltation) — **analytical ability** and **precision** reach **highest level**. You'll be **rigorous thinking**, **detail-oriented**, **highly efficient**. Excellent time for **optimizing work**, **health management**. Energy is **precise and efficient**; suitable for **analysis, organizing, service**.",
		"venus_exaltation": "Venus in **Pisces** (exaltation) — **love and compassion** are in **most ideal state**. You'll be **compassionate**, **artistically sensitive**, **emotionally rich**. Good period for **artistic creation**, **emotional connection**. Energy is **gentle and dreamy**; suitable for **creation, meditation, helping others**.",
		"mars_exaltation": "Mars in **Capricorn** (exaltation) — **action** and **endurance** reach **strongest state**. You'll be **clear goals**, **strong execution**, **persistent**. Excellent time for **building career**, **overcoming difficulties**. Energy is **firm and lasting**; suitable for **work, competition, pursuing goals**.",
		"jupiter_exaltation": "Jupiter in **Cancer** (exaltation) — **expansion** and **nurturing** are in **most ideal state**. You'll have **happy family**, **rich emotions**, **increased opportunities**. Good period for **improving family**, **emotional growth**. Energy is **nurturing and expanding**; suitable for **family activities, emotional exchange, building security**.",
		"saturn_exaltation": "Saturn in **Libra** (exaltation) — **responsibility** and **balance** reach **most harmonious state**. You'll be **mature relationships**, **fair and just**, **clear structure**. Excellent time for **improving relationships**, **building cooperation**. Energy is **balanced and stable**; suitable for **cooperation, negotiation, building structure**.",
	}
	if text, ok := interpretations[key]; ok {
		return text
	}
	planetName := map[string]string{
		"sun": "Sun", "moon": "Moon", "mercury": "Mercury", "venus": "Venus", "mars": "Mars",
		"jupiter": "Jupiter", "saturn": "Saturn", "uranus": "Uranus", "neptune": "Neptune", "pluto": "Pluto",
	}[string(planet)]
	dignityName := map[string]string{
		"domicile": "domicile", "exaltation": "exaltation", "detriment": "detriment", "fall": "fall",
	}[dignityType]
	return planetName + " is in **" + dignityName + "** state; energy expression is **" + (map[string]string{"domicile": "natural and smooth", "exaltation": "enhanced and elevated", "detriment": "restricted", "fall": "needs adjustment"}[dignityType]) + "**. Notice **using advantages**, **compensating for weaknesses**."
}

func getRussianDignityInterpretation(planet models.PlanetID, dignityType string) string {
	key := string(planet) + "_" + dignityType
	interpretations := map[string]string{
		"sun_domicile":       "Солнце в **Льве** (дом) — **самовыражение** и **лидерство** достигают **естественного пика**. Вы будете **уверены**, **лучисты**, **привлекать внимание**. Отличное время для **проявления таланта**, **ролей лидера**. Энергия **стабильна и мощна**; подходит для **творчества, выступлений, лидерства в проектах**.",
		"moon_domicile":      "Луна в **Раке** (дом) — **эмоции** и **интуиция** в **самом комфортном состоянии**. Вы будете **эмоционально богаты**, **защитны**, **сильное семейное сознание**. Хороший период для **обработки семейных дел**, **эмоционального исцеления**. Энергия **питательна и безопасна**; подходит для **заботы о других, домашнего декора, эмоционального выражения**.",
		"mercury_domicile":   "Меркурий в **Близнецах** или **Деве** (дом) — **мышление** и **общение** достигают **лучшего состояния**. Вы будете **сообразительны**, **ясны в выражениях**, с **высокой способностью к обучению**. Отличное время для **учебы, письма, переговоров**. Энергия **гибкая и точная**; подходит для **преподавания, анализа, обработки информации**.",
		"venus_domicile":     "Венера в **Тельце** или **Весах** (дом) — **любовь и красота** в **самом гармоничном состоянии**. Вы будете **обаятельны**, **эстетически восприимчивы**, **отношения гармоничны**. Хороший период для **углубления связей**, **художественного творчества**. Энергия **элегантная и сбалансированная**; подходит для **свиданий, искусства, наслаждения жизнью**.",
		"mars_domicile":      "Марс в **Овне** или **Скорпионе** (дом) — **действенность** и **сила воли** достигают **максимума**. Вы будете **энергичны**, **целеустремленны**, с **высокой исполнительностью**. Отличное время для **запуска проектов**, **преодоления препятствий**. Энергия **прямая и мощная**; подходит для **соревнований, спорта, достижения целей**.",
		"jupiter_domicile":   "Юпитер в **Стрельце** или **Рыбах** (дом) — **расширение** и **мудрость** в **самом идеальном состоянии**. У вас будет **широкое видение**, **оптимизм**, **больше возможностей**. Хороший период для **обучения, путешествий, поиска идеалов**. Энергия **щедрая и мудрая**; подходит для **образования, философии, исследований**.",
		"saturn_domicile":    "Сатурн в **Козероге** или **Водолее** (дом) — **ответственность** и **структура** достигают **самого стабильного состояния**. Вы станете **практичнее**, **дисциплинированнее**, с **четким долгосрочным планированием**. Отличное время для **построения карьеры**, **принятия ответственности**. Энергия **стабильная и прочная**; подходит для **работы, планирования, укрепления авторитета**.",
		"sun_exaltation":     "Солнце в **Овне** (экзальтация) — **пионерский дух** и **действенность** **значительно усилены**. Вы будете **полны жизненных сил**, **готовы к риску**, с **выдающимся лидерством**. Хорошее время для **запуска новых проектов**, **проявления смелости**. Энергия **активная и импульсивная**; подходит для **соревнований, спорта, лидерства**.",
		"moon_exaltation":    "Луна в **Тельце** (экзальтация) — **эмоции** и **безопасность** в **самом комфортном состоянии**. Вы будете **эмоционально стабильны**, **наслаждаться жизнью**, с **материальным удовлетворением**. Хороший период для **создания безопасности**, **наслаждения жизнью**. Энергия **питательная и стабильная**; подходит для **отдыха, еды, искусства**.",
		"mercury_exaltation": "Меркурий в **Деве** (экзальтация) — **аналитические способности** и **точность** достигают **высшего уровня**. У вас будет **строгое мышление**, **внимание к деталям**, **высокая эффективность**. Отличное время для **оптимизации работы**, **управления здоровьем**. Энергия **точная и эффективная**; подходит для **анализа, организации, служения**.",
		"venus_exaltation":   "Венера в **Рыбах** (экзальтация) — **любовь и сострадание** в **самом идеальном состоянии**. Вы будете **сострадательны**, **эстетически чувствительны**, **эмоционально богаты**. Хороший период для **художественного творчества**, **эмоциональной связи**. Энергия **нежная и мечтательная**; подходит для **творчества, медитации, помощи другим**.",
		"mars_exaltation":    "Марс в **Козероге** (экзальтация) — **действенность** и **выносливость** достигают **максимума**. У вас будут **четкие цели**, **высокая исполнительность**, **настойчивость**. Отличное время для **построения карьеры**, **преодоления трудностей**. Энергия **твердая и прочная**; подходит для **работы, соревнований, достижения целей**.",
		"jupiter_exaltation": "Юпитер в **Раке** (экзальтация) — **расширение** и **питание** в **самом идеальном состоянии**. У вас будет **счастливая семья**, **богатые эмоции**, **больше возможностей**. Хороший период для **улучшения семьи**, **эмоционального роста**. Энергия **питательная и расширяющая**; подходит для **семейных дел, эмоционального обмена, создания безопасности**.",
		"saturn_exaltation":  "Сатурн в **Весах** (экзальтация) — **ответственность** и **баланс** достигают **самого гармоничного состояния**. У вас будут **зрелые отношения**, **справедливость**, **четкая структура**. Отличное время для **улучшения отношений**, **создания сотрудничества**. Энергия **сбалансированная и стабильная**; подходит для **сотрудничества, переговоров, построения структуры**.",
		"sun_detriment":      "Солнце в **Водолее** (изгнание) — **самовыражение** может быть **ограничено**. Вы станете **независимее**, с **сильным чувством отстраненности**, **трудностями в проявлении**. Это время для **переопределения себя**, **погони за идеалами**, но следите за тем, чтобы **не быть слишком отстраненным**. Энергия **уникальная и обособленная**, нужен **баланс личного и коллективного**.",
		"moon_detriment":     "Луна в **Козероге** (изгнание) — **эмоции** могут быть **подавлены**. Вы будете **сдержанны**, **серьезны**, с **трудностями в выражении**. Это время для **принятия ответственности**, **создания структуры**, но следите за тем, чтобы **не подавлять чувства слишком сильно**. Энергия **серьезная и сдержанная**, нужно **позволить эмоциям течь**.",
		"mercury_detriment":  "Меркурий в **Стрельце** или **Рыбах** (изгнание) — **мышление** может быть **недостаточно точным**. Ваши **мысли будут перескакивать**, **идеализироваться**, **легко рассеиваться**. Это время для **погони за идеалами**, **интуитивного мышления**, но следите за тем, чтобы **не быть слишком идеалистичным**. Энергия **идеалистичная и рассеянная**, нужно **сохранять сосредоточенность**.",
		"venus_detriment":    "Венера в **Овне** или **Скорпионе** (изгнание) — **любовь и красота** могут **выражаться с трудом**. Ваши **отношения будут прямыми**, **эмоции сильными**, **трудно достичь баланса**. Это время для **погони за страстью**, **глубокой связи**, но следите за тем, чтобы **не впадать в крайности**. Энергия **сильная и прямая**, нужны **баланс и гармония**.",
		"mars_detriment":     "Марс в **Тельце** или **Весах** (изгнание) — **действенность** может быть **затруднена**. Вы будете **действовать медленно**, **колебаться**, с **трудностями в запуске**. Это время для **укрепления основ**, **поиска баланса**, но следите за тем, чтобы **не быть слишком пассивным**. Энергия **стабильная и медленная**, нужно **активно продвигаться**.",
		"jupiter_detriment":  "Юпитер в **Близнецах** или **Деве** (изгнание) — **расширение** может быть **ограничено**. У вас будет **узкое видение**, **излишняя осторожность**, **меньше возможностей**. Это время для **изучения деталей**, **построения основ**, но следите за тем, чтобы **не быть слишком ограниченным**. Энергия **осторожная и локальная**, нужно **расширять кругозор**.",
		"saturn_detriment":   "Сатурн в **Раке** или **Льве** (изгнание) — **ответственность** может **выражаться с трудом**. У вас будет **тяжелое эмоциональное бремя**, **трудности в создании структуры**, **слабое чувство авторитета**. Это время для **работы с чувствами**, **построения уверенности**, но следите за тем, чтобы **не перегружать себя**. Энергия **эмоциональная и обремененная**, нужно **выстраивать границы**.",
		"sun_fall":           "Солнце в **Весах** (падение) — **самовыражению** может потребоваться **баланс**. Вы будете **ценить отношения**, с **трудностями в независимости**, **необходимостью сотрудничества**. Это время для **улучшения отношений**, **поиска баланса**, но следите за тем, чтобы **не потерять себя**. Энергия **сбалансированная и зависимая**, нужно **сохранять индивидуальность**.",
		"moon_fall":          "Луна в **Скорпионе** (падение) — **эмоции** могут быть **слишком сильными**. Вы будете **глубоко чувствовать**, с **желанием контроля**, **трудностями в отпускании**. Это время для **глубокой трансформации**, **работы с тенью**, но следите за тем, чтобы **не контролировать слишком сильно**. Энергия **глубокая и интенсивная**, нужны **освобождение и доверие**.",
		"mercury_fall":       "Меркурий в **Рыбах** (падение) — **мышление** может быть **недостаточно ясным**. Ваша **интуиция усилится**, **легко запутаться**, **трудно быть точным**. Это время для **интуитивного мышления**, **художественного творчества**, но следите за тем, чтобы **не быть слишком туманным**. Энергия **интуитивная и размытая**, нужно **сохранять ясность**.",
		"venus_fall":         "Венера в **Деве** (падение) — **любовь и красота** могут **выражаться с трудом**. Вы будете **слишком придирчивы**, с **трудностями в получении удовольствия**, **напряженностью в отношениях**. Это время для **улучшения отношений**, **погони за совершенством**, но следите за тем, чтобы **не быть слишком критичным**. Энергия **придирчивая и напряженная**, нужно **принимать несовершенство**.",
		"mars_fall":          "Марс в **Раке** (падение) — на **действенность** могут **влиять эмоции**. Вы будете **движимы чувствами**, **действовать пассивно**, с **трудностями в прямоте**. Это время для **работы с эмоциями**, **создания безопасности**, но следите за тем, чтобы **не быть слишком эмоциональным**. Энергия **эмоциональная и пассивная**, нужны **рациональные действия**.",
		"jupiter_fall":       "Юпитер в **Козероге** (падение) — **расширение** может быть **ограничено**. Вы будете **излишне осторожны**, **возможностей меньше**, **видение сузится**. Это время для **построения основ**, **принятия ответственности**, но следите за тем, чтобы **не быть слишком ограниченным**. Энергия **осторожная и локальная**, нужно **расширять кругозор**.",
		"saturn_fall":        "Сатурн в **Овне** (падение) — **ответственность** может **выражаться с трудом**. Вы будете **действовать импульсивно**, с **трудностями в создании структуры**, **нехваткой терпения**. Это время для **быстрых действий**, **проявления смелости**, но следите за тем, чтобы **не быть слишком импульсивным**. Энергия **импульсивная и бесструктурная**, нужно **вырабатывать дисциплину**.",
	}
	if text, ok := interpretations[key]; ok {
		return text
	}
	planetName := map[string]string{
		"sun": "Солнце", "moon": "Луна", "mercury": "Меркурий", "venus": "Венера", "mars": "Марс",
		"jupiter": "Юпитер", "saturn": "Сатурн", "uranus": "Уран", "neptune": "Нептун", "pluto": "Плутон",
	}[string(planet)]
	dignityName := map[string]string{
		"domicile": "дом", "exaltation": "экзальтация", "detriment": "изгнание", "fall": "падение",
	}[dignityType]
	return planetName + " в состоянии **" + dignityName + "**; выражение энергии **" + (map[string]string{"domicile": "естественно и плавно", "exaltation": "усилено и возвышено", "detriment": "ограничено", "fall": "требует корректировки"}[dignityType]) + "**. Замечайте **использование преимуществ**, **компенсацию слабостей**."
}

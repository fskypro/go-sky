package fsxml_test

import (
	"fmt"
	"strings"

	"fsky.pro/fsserializer/fsxml"
)

const (
	_xml = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE junit SYSTEM "junit-result.dtd">
<root>
	<xml:fsky xmlns:a="https://fsky.pro" tests="2" failures="0" time="0.009" url="fsky.pro/fsserializer/fsxml">
		<properties>
			<property name="go.version">go1.13.1</property>
		</properties>
		<test classname="fsxml" id="ExampleParseXML" values="2.3 4.5"></test>

		<value ak='v100'> 100 </value>
		<value ak="v200"> 200 </value>
		<value ak="==[300]"> 300 </value>
		<value ak="400'500"> 400 </value>
		<value ak='600&quot;700&apos;800'> 
			<inner> xxx </inner>
		</value>

		<values> 100 200 300 </values>
		<items>
			<item> abcd </item>
			<item> efgh </item>
			<item><![CDATA[jkl>fa]]></item>
		</items>
		<中文:TAG>
			奥迪嘎达
		</中文:TAG>
	</xml:fsky>
</root>`
)

var _doc *fsxml.S_Doc

// 解释 XML
func ExampleParseXML() {
	var err error
	_doc, err = fsxml.LoadString(_xml)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// 获取测试
func ExampleGetting() {
	root := _doc.Root()            // 获得根节点
	node := root.Child("xml:fsky") // 获取指定名称的子节点

	fmt.Println("//", strings.Repeat("-", 50), "\n//", "获取：")
	fmt.Printf("doc.Header = %v\n", _doc.Header)
	fmt.Printf("fsky.root = %v\n", root == node.Root())               // 通过子孙节点获取根节点
	fmt.Printf("fsky.name = %v\n", node.Name())                       // 获取节点名称
	fmt.Printf("fsky.children.count = %v\n", node.ChildCount())       // 获取子节点个数
	fmt.Printf("fsky.LastChild.Name = %s\n", node.LastChild().Name()) // 获取最后一个子节点

	fmt.Printf("root.Child(\"xml:fsky/value\").Text() = %s\n",
		root.Child("xml:fsky/value").Text()) // 获取子孙节点
	fmt.Printf("root.Child(\"xml:fsky/value[1]\").Text() = %s\n",
		root.Child("xml:fsky/value[1]").Text()) // 获取指定路径的子孙节点，如果有同名的孙节点，则通过下标引用指定索引子节点
	fmt.Printf("root.Child(\"xml:fsky/value[-2]\").Text() = %s\n",
		root.Child("xml:fsky/value[-2]").Text()) // 获取指定路径的子孙节点，如果有同名的孙节点，则通过下标引用指定索引子节点，负索引表示后序（如 -1 表示最后一个）
	fmt.Printf("root.Child(\"xml:fsky/value[-2]\").Text() = %s\n",
		root.Child("xml:fsky/value[-2]").Text()) // 获取指定路径的子孙节点，如果有同名的孙节点，则通过下标引用指定索引子节点，负索引表示后序（如 -1 表示最后一个）
	fmt.Printf("root.Child(\"xml:fsky/value[ak=v200]\").Text() = %s\n",
		root.Child("xml:fsky/value[ak=v200]").Text()) // 获取指定路径的子孙节点，并且要求子节点的属性值与下标指定的一致
	fmt.Printf("root.Child(\"xml:fsky/value[ak='==[300]']\").Text() = %s\n",
		root.Child("xml:fsky/value[ak='==[300]']").Text()) // 获取指定路径的子孙节点，并且要求子节点的属性值与下标指定的一致，如果属性值中有中括号，则需要给属性值加上双引号或单引号，或·号
	fmt.Printf("root.Child(\"xml:fsky/value[ak='==[300]']\").Text() = %s\n",
		root.Child("xml:fsky/value[ak='==[300]']").Text()) // 获取指定路径的子孙节点，并且要求子节点的属性值与下标指定的一致，如果属性值中有中括号，则需要给属性值加上双引号或单引号，或·号
	fmt.Println(`root.Child("xml:fsky/value[ak='400"500']").Text() =`,
		root.Child(`xml:fsky/value[ak="400'500"]`).Text()) // 获取指定路径的子孙节点，并且要求子节点的属性值与下标指定的一致，属性值可以用三种引号括回：双引号、单引号、点括号
	fmt.Println("root.Child(\"xml:fsky/value[ak=`600\"700'800`]/inner\").Text() =",
		root.Child("xml:fsky/value[ak=`600\"700'800`]/inner").Text()) // 获取指定路径的子孙节点，并且要求子节点的属性值与下标指定的一致，属性值可以用三种引号括回：双引号、单引号、点括号

	node.Child("items").ChildByIndex(2).SetIsCData(true) // 获取指定索引的子节点

	fmt.Printf("fsky.attributes.count = %v\n", node.AttrCount())                                     // 获取节点的属性个数
	fmt.Printf("fsky.time = %f\n", node.Attr("time").AsFloat32(0.0))                                 // 以指定类型获取节点属性的值
	fmt.Printf("fsky.xmlns:a = %s:%s\n", node.Attr("xmlns:a").Space(), node.Attr("xmlns:a").Local()) // 获取属性名称的命名空间和名字部分
	fmt.Printf("fsky.test.id = %q\n", node.Child("test").ID())                                       // 获取节点的 ID 属性值（如果节点没有ID属性，则返回空字符串）
	fmt.Printf("fsky.test.[values] = %v\n", node.Child("test").Attr("values").AsFloat32s(nil))       // 以浮点 slice 形式，获取属性值

	fmt.Printf("fsky.value.text.toInt16 = %d\n", node.Child("value").AsInt16(0))       // 以指定类型获取节点内容
	fmt.Printf("fsky.values.text.toInt32s = %v\n", node.Child("values").AsInt32s(nil)) // 以指定类型或一组同名称的子节点内容
	fmt.Printf("fsky.items.[item] = %v\n", node.Child("items").ReadItem("item"))       // 获取指定名称子节点的内容（相当于：node.Child("items/item").Text()）

	fmt.Printf("xml.<namespace> = %s\n", _doc.GetNamespace("xml")) // 获取指定名称的名字空间
	fmt.Printf("a.<namespace> = %s\n", _doc.GetNamespace("a"))     // 来自于 <xml:fsky xmlns:a=...>
}

// 设置测试
func ExampleSetting() {
	root := _doc.Root()            // 获得根节点
	node := root.Child("xml:fsky") // 获取指定名称的子节点
	fmt.Println("\n//", strings.Repeat("-", 50), "\n//", "设置：")
	sub := fsxml.CreateNode(root.Doc(), "12tagl", "new add node") // 创建一个新节点
	if sub == nil {
		fmt.Println("错误的 xml tag 名称")
	}
	sub = fsxml.CreateNode(node.Doc(), "xml:tag", "new add node")
	node.AddChild(sub)
	attr := fsxml.NewAttr("attr name", "12 34 56") // 新建一个属性
	if attr == nil {
		fmt.Println("错误的属性名称")
	}
	attr = fsxml.NewAttr("xml:attr", "12 34 56")
	sub.SetAttr(attr)
}

// 保存测试
func ExampleSaving() {
	// 保存为非格式化 xml 文件
	fmt.Println("\n//", strings.Repeat("-", 50), "\n//", "保存为 xml 文件：")
	err := _doc.Save("./test.xml") // 保存为非格式化的 xml 文件
	if err != nil {
		panic(fmt.Sprintf("save xml file error, err: %s\n", err.Error()))
	} else {
		fmt.Println("./test.xml 文件保存成功")
	}

	// 保存为格式化 xml 文件
	doc := _doc.Clone()                                       // 复制一个文档
	err = doc.SaveIndent("./test_indent.xml", "\t", fsxml.LF) // 保存为格式化的 xml 文件
	if err != nil {
		panic(fmt.Sprintf("save indent xml file error, err: %s\n", err.Error()))
	} else {
		fmt.Println("./test_indent.xml 文件保存成功")
	}
}

// 加载 xml 文件测试
func ExampleLoading() {
	// 读取 xml 文件
	fmt.Println("\n//", strings.Repeat("-", 50), "\n//", "读取 xml 文件：")
	doc, err := fsxml.LoadFile("./test_indent.xml")
	if err != nil {
		panic(fmt.Sprintf("load indent xml file fail, error: %s\n", err.Error()))
	}
	fmt.Printf("doc's children count = %d\n", doc.Root().ChildCount())
}

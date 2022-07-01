export default function Dog(props) {
  const { type, position } = { ...props };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        gridColumn: position.y + 1, // css grids start from index 1 :(
        gridRow: position.x + 1,
      }}
    >
      <img height="75" width="75" src={`dogs/${type}.png`} />
      <span>{`{x: ${position.x}, y: ${position.y}}`}</span>
    </div>
  );
}
